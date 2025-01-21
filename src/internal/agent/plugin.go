package agent

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/mwantia/prometheus/internal/plugin/essentials"
	"github.com/mwantia/prometheus/internal/registry"
	"github.com/mwantia/prometheus/pkg/plugin"
	"github.com/mwantia/prometheus/pkg/plugin/cache"
	"github.com/mwantia/prometheus/pkg/plugin/identity"
	"github.com/mwantia/prometheus/pkg/plugin/tools"
)

var Plugins = map[string]PluginServe{
	"essentials": func() error {
		return essentials.Serve()
	},
}

type PluginServe func() error

func (a *PrometheusAgent) serveLocalPlugins() error {
	files, err := os.ReadDir(a.Config.PluginDir)
	if err != nil {
		return fmt.Errorf("unable to read directory '%s': %v", a.Config.PluginDir, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			path := fmt.Sprintf("%s/%s", a.Config.PluginDir, file.Name())
			if err := a.RunLocalPlugin(path); err != nil {
				log.Printf("Unable to load local plugin: %v", err)
			}
		}
	}

	return nil
}

func (a *PrometheusAgent) serveEmbedPlugins() error {
	for _, name := range a.Config.EmbedPlugins {
		p, exists := Plugins[name]
		if exists && p != nil {
			err := a.RunEmbedPlugin(name)
			if err != nil {
				log.Printf("Error serving embed plugin: %v", err)
			}
		} else {
			log.Printf("Embedded plugin '%s' doesn't exist", name)
		}
	}

	return nil
}

func (a *PrometheusAgent) RunEmbedPlugin(name string) error {
	path, err := os.Executable()
	if err != nil {
		return nil
	}

	if err = a.RunLocalPlugin(path, "plugin", name); err != nil {
		return fmt.Errorf("unable to load embbeded plugin '%s': %v", name, err)
	}

	return nil
}

func (a *PrometheusAgent) RunLocalPlugin(path string, arg ...string) error {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	log.Printf("Run local plugin '%s' with args: %v", path, arg)

	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig: plugin.Handshake,
		Plugins:         plugin.PluginMap,
		Cmd:             exec.Command(path, arg...),
	})

	rpc, err := client.Client()
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to create rpc-client: %w", err)
	}

	raw, err := rpc.Dispense("identity")
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to dispense 'identity' from plugin: %w", err)
	}

	ident := raw.(identity.IdentityService)
	info, err := ident.GetPluginInfo()
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to get plugin info: %w", err)
	}

	/* data, err := a.Config.GetPluginConfigMap(info.Name)
	if err != nil {
		log.Printf("Unable to load plugin config: %v", err)
	} */

	i := &registry.PluginInfo{
		Name:     info.Name,
		Version:  info.Version,
		Author:   info.Author,
		Metadata: info.Metadata,

		Services: registry.PluginServices{
			Identity: ident,
		},
		Cleanup: func() error {
			client.Kill()
			return nil
		},
	}

	for index, svr := range info.Services {
		switch svr.Type {
		case identity.ToolServiceType:
			key := fmt.Sprintf("tool.%v", index)
			raw, err := rpc.Dispense(key)
			if err != nil {
				return fmt.Errorf("failed to dispense service: %w", err)
			}

			service, success := raw.(tools.ToolService)
			if !success {
				return fmt.Errorf("failed to cast service")
			}

			i.Services.Tools = append(i.Services.Tools, service)

		case identity.CacheServiceType:
			raw, err := rpc.Dispense("cache")
			if err != nil {
				return fmt.Errorf("failed to dispense service: %w", err)
			}

			service, success := raw.(cache.CacheService)
			if !success {
				return fmt.Errorf("failed to cast service")
			}

			service.SetCache(&cache.SetCacheRequest{
				Key:   "foo",
				Value: []byte("bar"),
			})
		}
	}

	if err := a.Registry.RegisterPlugin(i); err != nil {
		return fmt.Errorf("failed to register plugin: %w", err)
	}

	log.Printf("Loaded local plugin '%v'", info.Name)

	/*if err := driver.Setup(plugin.PluginSetup{
		Data: data,
	}); err != nil {
		client.Kill()
		return err
	}*/

	return nil
}

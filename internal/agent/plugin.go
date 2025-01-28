package agent

import (
	"fmt"
	"os"
	"os/exec"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/mwantia/queueverse/internal/plugin/essentials"
	"github.com/mwantia/queueverse/pkg/plugin"
)

var Plugins = map[string]PluginServe{
	"essentials": func() error {
		return essentials.Serve()
	},
}

type PluginServe func() error

func (a *Agent) serveLocalPlugins() error {
	files, err := os.ReadDir(a.Config.PluginDir)
	if err != nil {
		return fmt.Errorf("unable to read directory '%s': %v", a.Config.PluginDir, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			path := fmt.Sprintf("%s/%s", a.Config.PluginDir, file.Name())
			if err := a.RunLocalPlugin(path); err != nil {
				a.Log.Warn("Unable to load local plugin", "path", path, "error", err)
			}
		}
	}

	return nil
}

func (a *Agent) serveEmbedPlugins() error {
	for _, name := range a.Config.EmbedPlugins {
		p, exists := Plugins[name]
		if exists && p != nil {
			err := a.RunEmbedPlugin(name)
			if err != nil {
				a.Log.Error("Error serving embed plugin", "name", name, "error", err)
			}
		} else {
			a.Log.Warn("Embedded plugin doesn't exist", "name", name)
		}
	}

	return nil
}

func (a *Agent) RunEmbedPlugin(name string) error {
	path, err := os.Executable()
	if err != nil {
		return nil
	}

	if err = a.RunLocalPlugin(path, "plugin", name); err != nil {
		return fmt.Errorf("unable to load embbeded plugin '%s': %v", name, err)
	}

	return nil
}

func (a *Agent) RunLocalPlugin(path string, arg ...string) error {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	a.Log.Debug("Running local plugin", "path", path, "args", arg)

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

	r, err := rpc.Dispense("plugin")
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to dispense 'plugin' from plugin: %w", err)
	}
	p := r.(plugin.Plugin)
	ts, err := p.GetTools()
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to dispense 'plugin' from plugin: %w", err)
	}

	a.Log.Debug("Tools", "map", ts)

	/* raw, err := rpc.Dispense("identity")
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

	metrics.RegisterActivePlugin(info.Name, info.Version, info.Author)
	for _, service := range info.Services {
		metrics.RegisterActiveService(info.Name, service.Name, string(service.Type))
	}

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

	a.Log.Info("Loaded new local plugin", "name", info.Name, "version", info.Version, "author", info.Author)
	*/
	return nil
}

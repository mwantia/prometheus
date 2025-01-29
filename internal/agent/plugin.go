package agent

import (
	"fmt"
	"os"
	"os/exec"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/mwantia/queueverse/internal/metrics"
	"github.com/mwantia/queueverse/internal/plugin/essentials"
	"github.com/mwantia/queueverse/internal/plugin/ollama"
	"github.com/mwantia/queueverse/pkg/plugin"
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/pkg/plugin/tools"
)

var Plugins = map[string]PluginServe{
	"essentials": func() error {
		return essentials.Serve()
	},
	"ollama": func() error {
		return ollama.Serve()
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
		HandshakeConfig: base.Handshake,
		Plugins:         plugin.Plugins,
		Cmd:             exec.Command(path, arg...),
		Logger:          a.Log.Named("plugin").Impl(),
	})

	rpc, err := client.Client()
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to create rpc connection: %w", err)
	}

	raw, err := rpc.Dispense(base.PluginBaseType)
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to dispense from plugin: %w", err)
	}

	basePlugin, exist := raw.(base.BasePlugin)
	if !exist {
		client.Kill()
		return fmt.Errorf("unable to cast raw interface")
	}

	info, err := basePlugin.GetPluginInfo()
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to get plugin info: %w", err)
	}

	cfgmap, err := a.Config.GetPluginConfigMap(info.Name)
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to load plugin config: %w", err)
	}

	if err := basePlugin.SetConfig(&base.PluginConfig{ConfigMap: cfgmap}); err != nil {
		client.Kill()
		return fmt.Errorf("failed to set plugin config: %w", err)
	}

	metrics.RegisterActivePlugin(info.Name, info.Version, "unknown")

	switch info.Type {
	case base.PluginProviderType:

		r, err := rpc.Dispense(base.PluginProviderType)
		if err != nil {
			client.Kill()
			return fmt.Errorf("failed to create rpc connection: %w", err)
		}

		providerPlugin, exist := r.(provider.ProviderPlugin)
		if !exist {
			client.Kill()
			return fmt.Errorf("unable to cast raw interface")
		}

		if err := providerPlugin.ProbePlugin(); err != nil {
			return fmt.Errorf("failed to probe plugin state: %w", err)
		}

		resp, err := providerPlugin.Chat(provider.ProviderChatRequest{
			Model: "llama3.2:latest",
			Messages: []provider.ProviderChatMessage{
				{
					Role:    "user",
					Content: "Why is the sky blue. Reply with 50 words or less.",
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to chat provider: %w", err)
		}

		a.Log.Warn("Provider chat response...", "model", resp.Model, "content", resp.Message.Content)

	case base.PluginToolsType:

		r, err := rpc.Dispense(base.PluginToolsType)
		if err != nil {
			client.Kill()
			return fmt.Errorf("failed to create rpc connection: %w", err)
		}

		toolPlugin := r.(tools.ToolPlugin)
		if err := toolPlugin.ProbePlugin(); err != nil {
			return fmt.Errorf("unable to probe plugin state: %w", err)
		}

	default:

		client.Kill()
		return fmt.Errorf("unknown plugin type '%s' is not supported", info.Type)
	}

	/*
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

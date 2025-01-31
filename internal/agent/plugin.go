package agent

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	hclog "github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/mwantia/queueverse/internal/metrics"
	"github.com/mwantia/queueverse/pkg/plugin"
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/pkg/plugin/tools"
	"github.com/mwantia/queueverse/plugins/mock"
	"github.com/mwantia/queueverse/plugins/ollama"
)

var Plugins = map[string]PluginServe{
	"mock": func() error {
		return plugin.ServeContext(func(ctx context.Context, logger hclog.Logger) interface{} {
			return &mock.MockProvider{
				Context: ctx,
				Logger:  logger,
			}
		})
	},
	"ollama": func() error {
		return plugin.ServeContext(func(ctx context.Context, logger hclog.Logger) interface{} {
			return &ollama.OllamaProvider{
				Context: ctx,
				Logger:  logger,
			}
		})
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

	plugin, exist := raw.(base.BasePlugin)
	if !exist {
		client.Kill()
		return fmt.Errorf("unable to cast raw interface")
	}

	info, err := plugin.GetPluginInfo()
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to get plugin info: %w", err)
	}

	cfgmap, err := a.Config.GetPluginConfigMap(info.Name)
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to load plugin config: %w", err)
	}

	if err := plugin.SetConfig(&base.PluginConfig{ConfigMap: cfgmap}); err != nil {
		client.Kill()
		return fmt.Errorf("failed to set plugin config: %w", err)
	}

	metrics.RegisterActivePlugin(info.Name, info.Version, info.Author)

	switch info.Type {
	case base.PluginProviderType:

		r, err := rpc.Dispense(base.PluginProviderType)
		if err != nil {
			client.Kill()
			return fmt.Errorf("failed to create rpc connection: %w", err)
		}

		plugin, exist := r.(provider.ProviderPlugin)
		if !exist {
			client.Kill()
			return fmt.Errorf("unable to cast raw interface")
		}

		if err := a.Registry.Register(info, plugin, func() error {
			client.Kill()
			return nil
		}); err != nil {
			client.Kill()
			return fmt.Errorf("failed to register plugin: %w", err)
		}

	case base.PluginToolsType:

		r, err := rpc.Dispense(base.PluginToolsType)
		if err != nil {
			client.Kill()
			return fmt.Errorf("failed to create rpc connection: %w", err)
		}

		plugin, exist := r.(tools.ToolPlugin)
		if !exist {
			client.Kill()
			return fmt.Errorf("unable to cast raw interface")
		}

		if err := a.Registry.Register(info, plugin, func() error {
			client.Kill()
			return nil
		}); err != nil {
			client.Kill()
			return fmt.Errorf("failed to register plugin: %w", err)
		}

	default:

		client.Kill()
		return fmt.Errorf("unknown plugin type '%s' is not supported", info.Type)
	}

	a.Log.Info("Loaded new local plugin", "name", info.Name, "version", info.Version, "author", info.Author)

	return nil
}

package agent

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/mwantia/prometheus/internal/plugin/debug"
	"github.com/mwantia/prometheus/internal/registry"
	"github.com/mwantia/prometheus/pkg/msg"
	"github.com/mwantia/prometheus/pkg/plugin"
)

var Plugins = map[string]plugin.Plugin{
	"debug": debug.NewPlugin(),
}

func (a *PrometheusAgent) serveLocalPlugins() error {
	files, err := os.ReadDir(a.Config.Agent.PluginDir)
	if err != nil {
		return fmt.Errorf("unable to read directory '%s': %v", a.Config.Agent.PluginDir, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			path := fmt.Sprintf("%s/%s", a.Config.Agent.PluginDir, file.Name())
			if err := a.RunLocalPlugin(path); err != nil {
				log.Printf("Unable to load local plugin: %v", err)
			}
		}
	}

	return nil
}

func (a *PrometheusAgent) serveEmbedPlugins() error {
	for _, name := range a.Config.Agent.EmbedPlugins {
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

	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig: plugin.Handshake,
		Plugins:         plugin.Plugins,
		Cmd:             exec.Command(path, arg...),
	})

	rpc, err := client.Client()
	if err != nil {
		client.Kill()
		return err
	}

	raw, err := rpc.Dispense("driver")
	if err != nil {
		client.Kill()
		return err
	}

	driver := raw.(plugin.Plugin)

	name, err := driver.Name()
	if err != nil {
		return err
	}

	cap, err := driver.GetCapabilities()
	if err != nil {
		return err
	}

	cfg := a.Config.GetPluginConfig(name)
	if cfg.Enabled {
		data, err := a.Config.GetPluginConfigMap(name)
		if err != nil {
			log.Printf("Unable to load plugin config: %v", err)
		}

		log.Printf("Loaded local plugin '%v'", name)

		s := plugin.PluginSetup{
			Data: data,
		}

		if a.Config.Agent.Kafka != nil {
			if s.Hub != nil {
				return fmt.Errorf("message hub already declared by another config")
			}

			s.Hub = &msg.KafkaMessageHub{
				Network:   a.Config.Agent.Kafka.Network,
				Address:   a.Config.Agent.Kafka.Address,
				Topic:     a.Config.Agent.Kafka.Topics,
				Partition: a.Config.Agent.Kafka.Partition,
			}
		}

		if err := driver.Setup(s); err != nil {
			client.Kill()
			return err
		}

		info := &registry.PluginInfo{
			Name:         name,
			Plugin:       driver,
			Capabilities: cap,

			Cleanup: func() error {
				client.Kill()
				return nil
			},
		}
		if err := a.Registry.RegisterPlugin(info); err != nil {
			return err
		}
	} else {
		log.Printf("Plugin '%s' is marked as disabled and will be ignored", name)
	}

	return nil
}

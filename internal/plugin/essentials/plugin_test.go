package essentials

import (
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/mwantia/queueverse/pkg/plugin"
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
)

func TestPlugin(t *testing.T) {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig: base.Handshake,
		Plugins:         plugin.Plugins,
		Cmd:             exec.Command("../../../build/queueverse", "plugin", "essentials"),
		Logger:          logger,
	})
	defer client.Kill()

	rpc, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	t.Run("Test.Plugin", func(t *testing.T) {
		raw, err := rpc.Dispense(base.PluginBaseType)
		if err != nil {
			log.Fatal(err)
		}

		b := raw.(base.BasePlugin)
		info, err := b.GetPluginInfo()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Type: %s", info.Type)
		log.Printf("Name: %s", info.Name)
		log.Printf("Version: %s", info.Version)

		switch info.Type {
		case base.PluginProviderType:
			raw, err := rpc.Dispense(base.PluginProviderType)
			if err != nil {
				log.Fatal(err)
			}

			prov := raw.(provider.ProviderPlugin)
			resp, err := prov.Chat(provider.ProviderChatRequest{
				Model: "test",
				Messages: []provider.ProviderChatMessage{
					{
						Role:    "user",
						Content: "This is a test",
					},
				},
			})
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Response: %s", resp.Message.Content)
		case base.PluginToolsType:

		default:
			log.Fatal("Unknown plugin type")
		}
	})
}

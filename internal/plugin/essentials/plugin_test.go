package essentials

import (
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/mwantia/queueverse/pkg/plugin"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
)

func TestPlugin(t *testing.T) {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig: plugin.Handshake,
		Plugins:         plugin.PluginMap,
		Cmd:             exec.Command("../../../build/queueverse", "plugin", "essentials"),
		Logger:          logger,
	})
	defer client.Kill()

	rpc, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	t.Run("Test.Plugin", func(t *testing.T) {
		raw, err := rpc.Dispense("plugin")
		if err != nil {
			log.Fatal(err)
		}

		plug := raw.(plugin.Plugin)
		prov, err := plug.GetProvider()
		if err != nil {
			log.Fatal(err)
		}

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

		log.Printf("Resp: %s", resp.Message.Content)
	})
}

package ollama

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/mwantia/queueverse/internal/config"
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
)

func TestOllama(t *testing.T) {
	cfg, err := config.ParseConfig("../../tests/config.hcl")
	if err != nil {
		t.Fatalf("Failed to parse test config: %v", err)
	}

	cfgmap, err := cfg.GetPluginConfigMap(PluginName)
	if err != nil {
		t.Fatalf("Failed to load plugin config: %v", err)
	}

	plugin := OllamaProvider{
		Context: context.TODO(),
	}
	if err := plugin.SetConfig(&base.PluginConfig{ConfigMap: cfgmap}); err != nil {
		t.Fatalf("Failed to set plugin config: %v", err)
	}

	t.Run("Ollama.Chat", func(t *testing.T) {
		resp, err := plugin.Chat(provider.ChatRequest{
			Model: "llama3.2:latest",
			Messages: []provider.ChatMessage{
				provider.UserMessage("Tell me the current time in germany."),
			},
			Tools: []provider.ToolDefinition{
				{
					Name:        "get_current_time",
					Description: "Get the current time in the specified timezone",
					Parameters: provider.ToolParameters{
						Type: provider.ToolTypeString,
						Required: []string{
							"timezone",
						},
						Properties: map[string]provider.ToolProperty{
							"timezone": {
								Type:        provider.ToolTypeString,
								Description: "The timezone to use. Must be a IANA Time Zone",
							},
						},
					},
				},
			},
		})
		if err != nil {
			t.Fatalf("Failed to perform chat request: %v", err)
		}

		debug, _ := json.Marshal(resp)
		log.Println(string(debug))
		t.Log(string(debug))
	})
}

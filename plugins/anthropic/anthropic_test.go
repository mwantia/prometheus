package anthropic

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/mwantia/queueverse/internal/config"
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
)

func TestAnthropic(t *testing.T) {
	cfg, err := config.ParseConfig("../../tests/config.hcl")
	if err != nil {
		t.Fatalf("Failed to parse test config: %v", err)
	}

	cfgmap, err := cfg.GetPluginConfigMap(PluginName)
	if err != nil {
		t.Fatalf("Failed to load plugin config: %v", err)
	}

	plugin := AnthropicProvider{
		Context: context.TODO(),
	}
	if err := plugin.SetConfig(&base.PluginConfig{ConfigMap: cfgmap}); err != nil {
		t.Fatalf("Failed to set plugin config: %v", err)
	}

	t.Run("Anthropic.Chat", func(t *testing.T) {
		resp, err := plugin.Chat(provider.ChatRequest{
			Model: string(anthropic.ModelClaude3Dot5HaikuLatest),
			Messages: []provider.ChatMessage{
				provider.UserMessage("Send a message to Roman Blake over Discord and tell him that I might arrive late to the meeting."),
			},
			Tools: []provider.ToolDefinition{
				{
					Name:        "send_discord_pm",
					Description: "Sends a private message over Discord to the specified user.",
					Parameters: provider.ToolParameters{
						Type:     provider.ToolTypeBoolean,
						Required: []string{},
						Properties: map[string]provider.ToolProperty{
							"username": {
								Type:        provider.ToolTypeString,
								Description: "The username the message will be send to",
							},
							"message": {
								Type:        provider.ToolTypeString,
								Description: "The message send to over Discord (Supports markdown).",
							},
						},
					},
				},
				{
					Name:        "get_discord_contact",
					Description: "Returns the username within Discord for the specified search.",
					Parameters: provider.ToolParameters{
						Type:     provider.ToolTypeString,
						Required: []string{},
						Properties: map[string]provider.ToolProperty{
							"search": {
								Type: provider.ToolTypeString,
								Description: `Defines the search query used to find the correct contact/username.
								This can be the displayname, surname, lastname or other available contact information`,
							},
						},
					},
				},
				/*{
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
				},*/
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

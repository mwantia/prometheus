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
		req := provider.ChatRequest{
			Model: string(anthropic.ModelClaude3Dot5HaikuLatest),
			Messages: []provider.ChatMessage{
				provider.UserMessage("Send a message to the username 'romanblake' over Discord and tell him that I might arrive late to the meeting."),
			},
			Tools: []provider.ToolDefinition{
				{
					Name: "send_discord_pm",
					Description: `Sends a private message over Discord to a user.
					The message property supports markdown to allow, for example bold or cursiv text.
					It should only be used when the user wants to send a message, or even a tool response to another user.
					You must define the intend of the message, as well as the user it was send or requested from.
					Use templates like '{{ .displayname }}' to refer to the user making the request (e.g. Displayname = 'Max Mustermann').`,
					Parameters: provider.ToolParameters{
						Type:     provider.ToolTypeBoolean,
						Required: []string{"user", "message"},
						Properties: map[string]provider.ToolProperty{
							"user": {
								Type:        provider.ToolTypeString,
								Description: `The user the message will be send to.`,
							},
							"message": {
								Type:        provider.ToolTypeString,
								Description: `The message send to over Discord (Supports markdown).`,
							},
						},
					},
				},
				{
					Name: "get_discord_contact",
					Description: `Retrieves a list of usernames available within Discord.
					Can be used in combination with other Discord tools that require information about a user.
					The received userdata is stored and provided to other tool calls as variables. 
					This can be accessed in a template format by defining '{{ user.<property> }}'.
					These can even be used in property values for other tool calls.
					The following variables will become available after searching for a user:
					* user.username
					* user.displayname
					* user.mail
					* user.status`,
					Parameters: provider.ToolParameters{
						Type:     provider.ToolTypeString,
						Required: []string{},
						Properties: map[string]provider.ToolProperty{
							"search": {
								Type: provider.ToolTypeString,
								Description: `Defines the search query used to find the correct user.
								This can be the displayname, surname, lastname or other available contact information`,
							},
						},
					},
				},
				{
					Name: "get_current_time",
					Description: `Get the current time in the specified timezone.
					The timezone must be a IANA compatible timezone.
					The output is in the following format 'Mon Jan 2 15:04:05'.
					Only use the toll, if the conversation specifically requires the current time.`,
					Parameters: provider.ToolParameters{
						Type:     provider.ToolTypeString,
						Required: []string{"timezone"},
						Properties: map[string]provider.ToolProperty{
							"timezone": {
								Type:        provider.ToolTypeString,
								Description: "The timezone to use. Must be a IANA Time Zone.",
							},
						},
					},
				},
			},
		}

		resp, err := plugin.Chat(req)
		if err != nil {
			t.Fatalf("Failed to perform chat request: %v", err)
		}

		debug, _ := json.Marshal(resp)
		log.Println(string(debug))

		req.Messages = append(req.Messages, resp.Messages...)
	})
}

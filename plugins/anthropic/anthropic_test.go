package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/mwantia/queueverse/internal/config"
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/plugins/anthropic/tools"
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
				// provider.UserMessage("Send a message to 'Roman Blake' over Discord and tell him that I might arrive late to the meeting."),
				provider.UserMessage("Tell me the current time in germany."),
			},
			Tools: []provider.ToolDefinition{
				tools.TimeGetCurrent,
				tools.DiscordListContact,
				tools.DiscordSendPM,
			},
		}

		resp, err := plugin.Chat(req)
		if err != nil {
			t.Fatalf("Failed to perform chat request: %v", err)
		}

		debug, _ := json.Marshal(resp)
		log.Println(string(debug))

		// This should result in a simple chat request without any additional tool calls
		if len(resp.Messages) == 1 && len(resp.Messages[0].ToolCalls) == 0 {
			log.Println(resp.Messages[0].Content)
		} else {
			for _, msg := range resp.Messages {
				if len(msg.ToolCalls) > 0 {
					for _, toolcall := range msg.ToolCalls {

						output, err := executeToolCall(toolcall)
						if err != nil {
							t.Fatalf("failed to execute tool call: %v", err)
						}

						req.Messages = append(req.Messages, provider.ChatMessage{
							ID:      msg.ID,
							Role:    provider.ChatRoleTool,
							Content: output,
						})
					}
				} else {
					req.Messages = append(req.Messages, msg)
				}
			}
		}

		resp, err = plugin.Chat(req)
		if err != nil {
			t.Fatalf("Failed to perform chat request: %v", err)
		}

		debug, _ = json.Marshal(resp)
		log.Println(string(debug))
	})
}

func executeToolCall(toolcall provider.ToolCall) (string, error) {
	switch toolcall.Function.Name {
	case tools.TimeGetCurrent.Name:
		timezone, exist := toolcall.Function.Arguments["timezone"]
		if !exist {
			return "", fmt.Errorf("failed too call '%s': argument 'timezone' not provided", toolcall.Function.Name)
		}

		location, err := time.LoadLocation(timezone.(string))
		if err != nil {
			return "", fmt.Errorf("failed to load location: ")
		}

		return time.Now().In(location).Format("Mon Jan 2 15:04:05"), nil

	case tools.DiscordListContact.Name:

	case tools.DiscordSendPM.Name:
	}
	return "", nil
}

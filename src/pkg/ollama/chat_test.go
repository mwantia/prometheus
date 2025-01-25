package ollama

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/mwantia/prometheus/internal/config"
)

var text strings.Builder

func TestChat(tst *testing.T) {
	cfg, err := config.ParseConfig("../../tests/config.hcl")
	if err != nil {
		tst.Fatalf("Unable to parse config: %v", err)
	}

	ctx := context.TODO()
	client := CreateClient(cfg.Ollama.Endpoint, "", http.DefaultClient)
	client.Style = ConciseStyle

	tools := []Tool{
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "get_current_time",
				Description: "Get the current time in the specified timezone",
				Parameters: ToolFunctionParameter{
					Type:     "string",
					Required: []string{"timezone"},
					Properties: map[string]ToolFunctionProperty{
						"timezone": {
							Type:        "string",
							Description: "The timezone to use. Must be a IANA Time Zone",
						},
					},
				},
			},
		},
	}

	tst.Run("Ollama.Chat", func(t *testing.T) {
		if err := client.ChatTools(ctx, ChatRequest{
			Model: cfg.Ollama.Model,
			Messages: []ChatMessage{
				{
					Role:    "user",
					Content: "Tell me the current time in germany.",
				},
			},
		}, handleChatResponse, handleToolCallResponse, tools); err != nil {
			tst.Errorf("Unable to perform chat request: %v", err)
		}

		log.Println(text.String())
		tst.Log(text.String())
	})
}

func handleChatResponse(resp ChatResponse) error {
	text.WriteString(resp.Message.Content)
	return nil
}

func handleToolCallResponse(tc ToolCall) (string, error) {
	switch tc.Function.Name {
	case "get_current_time":
		tz := tc.Function.Arguments["timezone"]
		loc, _ := time.LoadLocation(tz.(string))

		return time.Now().In(loc).Format("Mon Jan 2 15:04:05"), nil
	}

	return "", fmt.Errorf("no tool with the function name '%s' found.", tc.Function.Name)
}

package ollama

import (
	"context"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/mwantia/prometheus/internal/config"
)

func TestChat(tst *testing.T) {
	cfg, err := config.ParseConfig("../../tests/config.hcl")
	if err != nil {
		tst.Fatalf("Unable to parse config: %v", err)
	}

	ctx := context.TODO()
	client := CreateClient(cfg.Ollama.Address, http.DefaultClient)
	client.Style = ConciseStyle

	tst.Run("Ollama.Chat", func(t *testing.T) {
		var text strings.Builder

		if err := client.ChatTools(ctx, ChatRequest{
			Model: cfg.Ollama.Model,
			Messages: []ChatMessage{
				{
					Role:    "user",
					Content: "Tell me the current time in germany.",
				},
			},
		}, func(resp ChatResponse) error {
			text.WriteString(resp.Message.Content)
			return nil
		}, []Tool{
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
		}); err != nil {
			tst.Errorf("Unable to perform chat request: %v", err)
		}

		log.Println(text.String())
		tst.Log(text.String())
	})
}

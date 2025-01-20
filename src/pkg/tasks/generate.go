package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/mwantia/prometheus/internal/config"
	"github.com/mwantia/prometheus/pkg/ollama"
)

const TaskTypeGeneratePrompt = "task:generate"

type GeneratePrompt struct {
	Content string `json:"content"`
	Model   string `json:"model,omitempty"`
	Style   string `json:"style,omitempty"`
}

func HandleGeneratePromptTask(cfg *config.Config) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		client := ollama.CreateClient(cfg.Ollama.Address, http.DefaultClient)

		tools := createTools()
		var prompt GeneratePrompt
		if err := json.Unmarshal(t.Payload(), &prompt); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		if prompt.Model == "" {
			prompt.Model = cfg.Ollama.Model
		}

		log.Printf("Handling prompt: %s", prompt.Content)
		req := ollama.ChatRequest{
			Model: prompt.Model,
			Messages: []ollama.ChatMessage{
				{
					Role:    "user",
					Content: prompt.Content,
				},
			},
		}

		var text strings.Builder
		if err := client.ChatTools(ctx, req, func(resp ollama.ChatResponse) error {
			text.WriteString(resp.Message.Content)
			return nil
		}, tools); err != nil {
			return fmt.Errorf("failed chat tools: %w", err)
		}

		log.Println(text.String())
		if _, err := fmt.Fprint(t.ResultWriter(), text.String()); err != nil {
			return fmt.Errorf("failed to write task result: %w", err)
		}

		return nil
	}
}

func CreateGeneratePromptTask(content, model string) (*asynq.Task, error) {
	prompt, err := json.Marshal(GeneratePrompt{
		Content: content,
		Model:   model,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskTypeGeneratePrompt, prompt), nil
}

func createTools() []ollama.Tool {
	return []ollama.Tool{
		{
			Type: "function",
			Function: ollama.ToolFunction{
				Name:        "get_current_time",
				Description: "Get the current time in the specified timezone",
				Parameters: ollama.ToolFunctionParameter{
					Type:     "string",
					Required: []string{"timezone"},
					Properties: map[string]ollama.ToolFunctionProperty{
						"timezone": {
							Type:        "string",
							Description: "The timezone to use. Must be a IANA Time Zone",
						},
					},
				},
			},
		},
	}
}

package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/mwantia/prometheus/internal/config"
	"github.com/mwantia/prometheus/internal/registry"
	"github.com/mwantia/prometheus/pkg/ollama"
)

const TaskTypeGeneratePrompt = "task:generate"

type GeneratePrompt struct {
	Content string `json:"content"`
	Model   string `json:"model,omitempty"`
	Style   string `json:"style,omitempty"`
}

func CreateGeneratePromptTask(cfg *config.Config, reg *registry.PluginRegistry) func(context.Context, *asynq.Task) error {
	client := ollama.CreateClient(cfg.Ollama.Address, http.DefaultClient)
	tools := createTools()

	plugins := reg.GetPlugins()
	for _, plugin := range plugins {
		info, err := plugin.Services.Identity.GetPluginInfo()
		if err != nil {
			log.Printf("Unable to load plugin info for '%s'", plugin.Name)
		}

		for _, service := range info.Services {
			log.Printf("Name: %s", service.Name)
		}
	}

	return handleGeneratePromptTask(client, tools, cfg.Ollama.Model)
}

func handleGeneratePromptTask(client *ollama.Client, tools []ollama.Tool, model string) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var prompt GeneratePrompt
		if err := json.Unmarshal(t.Payload(), &prompt); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		if prompt.Model == "" {
			prompt.Model = model
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
		rHandler := func(r ollama.ChatResponse) error {
			text.WriteString(r.Message.Content)
			return nil
		}
		tHandler := func(tc ollama.ToolCall) (string, error) {
			log.Printf("Handling tool call: %s", tc.Function.Name)

			switch tc.Function.Name {
			case "get_current_time":
				tz := tc.Function.Arguments["timezone"]
				loc, _ := time.LoadLocation(tz.(string))

				return time.Now().In(loc).Format("Mon Jan 2 15:04:05"), nil
			}

			return "", fmt.Errorf("no tool with the function name '%s' found", tc.Function.Name)
		}

		if err := client.ChatTools(ctx, req, rHandler, tHandler, tools); err != nil {
			return fmt.Errorf("failed chat tools: %w", err)
		}

		log.Println(text.String())
		if _, err := fmt.Fprint(t.ResultWriter(), text.String()); err != nil {
			return fmt.Errorf("failed to write task result: %w", err)
		}

		return nil
	}
}

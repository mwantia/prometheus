package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/mwantia/queueverse/internal/config"
	"github.com/mwantia/queueverse/internal/registry"
	"github.com/mwantia/queueverse/pkg/log"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
)

func GenerateTaskId() string {
	return fmt.Sprintf("t%d", time.Now().UnixNano())
}

func CreateGenerateResponse(info *asynq.TaskInfo) (*GenerateResponse, error) {
	var result GenerateResult
	if len(info.Result) > 0 {
		if err := json.Unmarshal(info.Result, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal task response: %w", err)
		}
	}

	return &GenerateResponse{
		Task:   info.ID,
		State:  info.State.String(),
		Pool:   info.Queue,
		Result: result,
	}, nil
}

func CreateGenerateTask(req GenerateRequest) (*asynq.Task, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error during request marshal: %w", err)
	}

	return asynq.NewTask(TaskTypeGenerateName, payload), nil
}

func CreateGenerateTaskHandler(cfg *config.Config, registry *registry.Registry) func(context.Context, *asynq.Task) error {
	log := log.New("asynq")

	providers, _ := registry.GetProviders()

	return handleGenerateTask(log, providers)
}

func handleGenerateTask(log log.Logger, providers []provider.ProviderPlugin) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var req GenerateRequest

		startTime := time.Now()

		if err := json.Unmarshal(t.Payload(), &req); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		log.Info("Handling generate task", "model", req.Model, "content", req.Content)

		for _, prov := range providers {
			models, err := prov.GetModels()
			if err != nil {
				log.Warn("Unable to load models from provider", "error", err)
			}

			for _, model := range *models {
				if model.Name == req.Model {
					request := provider.ChatRequest{
						Model: req.Model,
						Messages: []provider.ChatMessage{
							provider.UserMessage(req.Content),
						},
					}

					resp, err := prov.Chat(request)
					if err != nil {
						log.Error("Failed to generate chat prompt")
					}

					duration := time.Since(startTime).Seconds()
					result := GenerateResult{
						Content: resp.Message.Content,
						Model:   resp.Model,
						Metadata: map[string]any{
							"duration": duration,
						},
					}

					// metrics.ClientGeneratePromptTasksDurationSeconds.WithLabelValues(oc.Endpoint, req.Model, req.Style).Observe(duration)

					log.Debug(resp.Message.Content)

					data, err := json.Marshal(result)
					if err != nil {
						return fmt.Errorf("failed to marshal final response: %w", err)
					}

					if _, err := t.ResultWriter().Write(data); err != nil {
						return fmt.Errorf("failed to write task result: %w", err)
					}
				}
			}
		}

		return nil
	}
}

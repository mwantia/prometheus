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

func CreateGenerateTaskHandler(cfg *config.Config, reg *registry.Registry) func(context.Context, *asynq.Task) error {
	return handleGenerateTask(log.New("asynq"), reg)
}

func handleGenerateTask(log log.Logger, reg *registry.Registry) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var req GenerateRequest

		startTime := time.Now()

		if err := json.Unmarshal(t.Payload(), &req); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		prov, err := reg.GetModelProvider(req.Model)
		if err != nil {
			return fmt.Errorf("failed to load model '%s': %w", req.Model, err)
		}

		log.Info("Handle Generate Task", "model", req.Model, "content", req.Content)

		request := provider.ChatRequest{
			Model: req.Model,
			Messages: []provider.ChatMessage{
				provider.UserMessage(req.Content),
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
		}

		resp, err := prov.Chat(request)
		if err != nil {
			return fmt.Errorf("failed to generate chat prompt: %w", err)
		}

		duration := time.Since(startTime).Seconds()
		result := GenerateResult{
			Content: resp.Messages[0].Content,
			Model:   resp.Model,
			Metadata: map[string]any{
				"duration": duration,
			},
		}

		// metrics.ClientGeneratePromptTasksDurationSeconds.WithLabelValues(oc.Endpoint, req.Model, req.Style).Observe(duration)

		debug, _ := json.Marshal(resp)
		log.Debug(string(debug))

		data, err := json.Marshal(result)
		if err != nil {
			return fmt.Errorf("failed to marshal final response: %w", err)
		}

		if _, err := t.ResultWriter().Write(data); err != nil {
			return fmt.Errorf("failed to write task result: %w", err)
		}

		return nil
	}
}

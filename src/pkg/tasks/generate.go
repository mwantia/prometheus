package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/mwantia/prometheus/internal/config"
	"github.com/mwantia/prometheus/internal/metrics"
	"github.com/mwantia/prometheus/internal/registry"
	"github.com/mwantia/prometheus/pkg/log"
	"github.com/mwantia/prometheus/pkg/ollama"
	"github.com/prometheus/client_golang/prometheus"
)

const TaskTypeGenerateName = "task:generate"

type GenerateRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model,omitempty"`
	Style  string `json:"style,omitempty"`
}

type GenerateResponse struct {
	Task   string `json:"task"`
	State  string `json:"state"`
	Pool   string `json:"pool"`
	Result string `json:"result,omitempty"`
}

func GenerateTaskId() string {
	return fmt.Sprintf("t%d", time.Now().UnixNano())
}

func CreateGenerateResponse(info *asynq.TaskInfo) GenerateResponse {
	return GenerateResponse{
		Task:   info.ID,
		State:  info.State.String(),
		Pool:   info.Queue,
		Result: string(info.Result),
	}
}

func CreateGenerateTask(req GenerateRequest) (*asynq.Task, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error during request marshal: %w", err)
	}

	return asynq.NewTask(TaskTypeGenerateName, payload), nil
}

func CreateGenerateTaskHandler(cfg *config.Config, reg *registry.PluginRegistry) func(context.Context, *asynq.Task) error {
	oc := ollama.CreateClient(cfg.Ollama.Endpoint, cfg.Ollama.Model, http.DefaultClient)
	ts := createTools()
	log := log.New("asynq")

	ps := reg.GetPlugins()
	for _, p := range ps {
		info, err := p.Services.Identity.GetPluginInfo()
		if err != nil {
			log.Error("Unable to load plugin info", "name", p.Name)
		}

		for _, s := range info.Services {
			log.Debug("Processing service...", "name", s.Name)
		}
	}

	return handleGenerateTask(log, oc, ts)
}

func handleGenerateTask(log *log.Logger, oc *ollama.Client, ts []ollama.Tool) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var req GenerateRequest

		if err := json.Unmarshal(t.Payload(), &req); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		if req.Model == "" {
			req.Model = oc.Model
		}
		if req.Style == "" {
			req.Style = string(oc.Style)
		}

		observer := prometheus.ObserverFunc(func(v float64) {
			metrics.ClientGeneratePromptTasksDurationSeconds.WithLabelValues(oc.Endpoint, req.Model, req.Style).Observe(v)
		})

		timer := prometheus.NewTimer(observer)
		defer timer.ObserveDuration()

		log.Info("Handling generate task...", "model", req.Model, "prompt", req.Prompt)

		creq := ollama.ChatRequest{
			Model: req.Model,
			Messages: []ollama.ChatMessage{
				{
					Role:    "user",
					Content: req.Prompt,
				},
			},
		}

		var text strings.Builder
		rhandler := func(r ollama.ChatResponse) error {
			text.WriteString(r.Message.Content)
			return nil
		}
		thandler := func(tc ollama.ToolCall) (string, error) {
			log.Info("Handling tool call...", "name", tc.Function.Name)

			switch tc.Function.Name {
			case "get_current_time":
				tz := tc.Function.Arguments["timezone"]
				loc, err := time.LoadLocation(tz.(string))
				if err != nil {
					return "", fmt.Errorf("unable to load timezone '%v': %w", tz, err)
				}

				return time.Now().In(loc).Format("Mon Jan 2 15:04:05"), nil
			}

			return "", fmt.Errorf("no tool with the function name '%s' found", tc.Function.Name)
		}

		if err := oc.ChatTools(ctx, creq, rhandler, thandler, ts); err != nil {
			return fmt.Errorf("failed chat tools: %w", err)
		}

		log.Debug(text.String())

		if _, err := fmt.Fprint(t.ResultWriter(), text.String()); err != nil {
			return fmt.Errorf("failed to write task result: %w", err)
		}

		return nil
	}
}

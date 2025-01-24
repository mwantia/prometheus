package ops

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/mwantia/prometheus/internal/config"
	"github.com/mwantia/prometheus/internal/registry"
	"github.com/mwantia/prometheus/pkg/log"
	"github.com/mwantia/prometheus/pkg/tasks"
)

type Client struct {
	Operation

	Log log.Logger
	mux *asynq.ServeMux
	srv *asynq.Server
}

func (c *Client) Create(cfg *config.Config, reg *registry.PluginRegistry) (Cleanup, error) {
	queues, err := validateAndConvertQueues(cfg.Client.Queues)
	if err != nil {
		return nil, fmt.Errorf("error validating or converting queues: %w", err)
	}

	c.Log = *log.New("client")
	c.mux = asynq.NewServeMux()
	c.srv = asynq.NewServer(asynq.RedisClientOpt{
		Addr:     cfg.Redis.Endpoint,
		DB:       cfg.Redis.Database,
		Password: cfg.Redis.Password,
	}, asynq.Config{
		Concurrency: 1,
		Queues:      queues,
	})

	c.mux.HandleFunc(tasks.TaskTypeGeneratePrompt, tasks.CreateGeneratePromptTask(cfg, reg))

	return func(ctx context.Context) error {
		c.srv.Shutdown()
		return nil
	}, nil
}

func (c *Client) Serve(ctx context.Context) error {
	return c.srv.Run(c.mux)
}

func validateAndConvertQueues(queues []string) (map[string]int, error) {
	valid := map[string]int{
		"debug":  10,
		"high":   7,
		"normal": 2,
		"low":    1,
	}

	result := make(map[string]int)
	for _, queue := range queues {
		if val, ok := valid[queue]; ok {
			result[queue] = val
		} else {
			return nil, fmt.Errorf("invalid queue: %s", queue)
		}
	}
	return result, nil
}

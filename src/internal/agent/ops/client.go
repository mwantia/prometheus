package ops

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/mwantia/prometheus/internal/config"
	"github.com/mwantia/prometheus/internal/registry"
	"github.com/mwantia/prometheus/pkg/tasks"
)

type Client struct {
	Operation

	mux *asynq.ServeMux
	srv *asynq.Server
}

func (c *Client) Create(cfg *config.Config, reg *registry.PluginRegistry) (Cleanup, error) {
	c.mux = asynq.NewServeMux()
	c.srv = asynq.NewServer(asynq.RedisClientOpt{
		Addr: cfg.Redis.Endpoint,
		DB:   cfg.Redis.Database,
	}, asynq.Config{
		Concurrency: 1,
		Queues: map[string]int{
			"critical": 7,
			"default":  2,
			"low":      1,
		},
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

package ops

import (
	"context"

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
	c.Log = *log.New("client")
	c.mux = asynq.NewServeMux()
	c.srv = asynq.NewServer(asynq.RedisClientOpt{
		Addr:     cfg.Redis.Endpoint,
		DB:       cfg.Redis.Database,
		Password: cfg.Redis.Password,
	}, asynq.Config{
		Concurrency: 1,
		Queues: map[string]int{
			cfg.PoolName: 10,
		},
	})

	c.mux.HandleFunc(tasks.TaskTypeGenerateName, tasks.CreateGenerateTaskHandler(cfg, reg))

	return func(ctx context.Context) error {
		c.srv.Shutdown()
		return nil
	}, nil
}

func (c *Client) Serve(ctx context.Context) error {
	return c.srv.Run(c.mux)
}

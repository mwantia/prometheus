package ops

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/mwantia/queueverse/internal/config"
	"github.com/mwantia/queueverse/internal/registry"
	"github.com/mwantia/queueverse/pkg/log"
	"github.com/mwantia/queueverse/pkg/tasks"
)

type Client struct {
	Operation

	Log log.Logger
	mux *asynq.ServeMux
	srv *asynq.Server
}

func (c *Client) Create(cfg *config.Config, registry *registry.Registry) (Cleanup, error) {
	c.Log = log.New("client")
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

	c.mux.HandleFunc(tasks.TaskTypeGenerateName, tasks.CreateGenerateTaskHandler(cfg, registry))

	return func(ctx context.Context) error {
		c.srv.Shutdown()
		return nil
	}, nil
}

func (c *Client) Serve(ctx context.Context) error {
	return c.srv.Run(c.mux)
}

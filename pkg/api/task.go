package api

import (
	"context"
	"net/http"
)

type Task interface {
	Queue(context.Context, string) error

	Wait(context.Context) error

	Result(context.Context) (string, error)
}

type taskImpl struct {
	client clientImpl
	config TaskConfig
	task   string
}

func NewTask(http *http.Client, cfg TaskConfig) Task {
	return &taskImpl{
		client: clientImpl{
			http: http,
			config: ClientConfig{
				Endpoint: cfg.Endpoint,
				Model:    cfg.Model,
				Token:    cfg.Token,
			},
		},
		config: cfg,
	}
}

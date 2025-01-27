package api

import (
	"context"
	"net/http"
)

type Client interface {
	QueueTask(context.Context, QueueRequest) (QueueResponse, error)

	QueueTaskResult(context.Context, string) (QueueResult, error)

	ListQueuedTasks(context.Context) ([]QueueResult, error)

	QueueState(context.Context, string) (bool, error)
}

type clientImpl struct {
	http   *http.Client
	config ClientConfig
}

func NewClient(http *http.Client, cfg ClientConfig) Client {
	return &clientImpl{
		http:   http,
		config: cfg,
	}
}

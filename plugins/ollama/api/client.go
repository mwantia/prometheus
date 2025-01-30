package api

import (
	"context"
	"net/http"
)

type Client interface {
	Health(context.Context) (bool, error)

	Tags(context.Context) ([]TagModel, error)

	Chat(context.Context, ChatRequest, ChatResponseHandler) error
}

type clientImpl struct {
	http   *http.Client
	config ClientConfig
}

type ClientConfig struct {
	Endpoint string `json:"endpoint"`
}

func CreateClient(http *http.Client, cfg ClientConfig) Client {
	return &clientImpl{
		http:   http,
		config: cfg,
	}
}

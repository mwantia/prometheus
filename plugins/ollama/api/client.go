package api

import (
	"context"
	"net/http"
)

type Client interface {
	Health(context.Context) (bool, error)

	Chat(context.Context, ChatRequest, ChatResponseHandler) error
}

type clientImpl struct {
	http   *http.Client
	config ClientConfig
}

type ClientConfig struct {
	Endpoint string `json:"endpoint"`
	Model    string `json:"model,omitempty"`
}

func CreateClient(http *http.Client, cfg ClientConfig) Client {
	return &clientImpl{
		http:   http,
		config: cfg,
	}
}

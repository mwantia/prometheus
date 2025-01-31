package main

import (
	"context"
	"log"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/mwantia/queueverse/pkg/plugin"
	"github.com/mwantia/queueverse/plugins/ollama"
)

func main() {
	if err := plugin.ServeContext(func(ctx context.Context, logger hclog.Logger) interface{} {
		return &ollama.OllamaProvider{
			Context: ctx,
			Logger:  logger,
		}
	}); err != nil {
		log.Panic(err)
	}
}

package ops

import (
	"context"

	"github.com/mwantia/queueverse/internal/config"
	"github.com/mwantia/queueverse/internal/registry"
)

type Operation interface {
	Create(*config.Config, *registry.PluginRegistry) (Cleanup, error)

	Serve(context.Context) error
}

type Cleanup func(context.Context) error

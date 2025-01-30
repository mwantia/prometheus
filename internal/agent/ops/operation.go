package ops

import (
	"context"

	"github.com/mwantia/queueverse/internal/config"
	"github.com/mwantia/queueverse/internal/registry"
)

type Operation interface {
	Create(*config.Config, *registry.Registry) (Cleanup, error)

	Serve(context.Context) error
}

type Cleanup func(context.Context) error

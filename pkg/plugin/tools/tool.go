package tools

import (
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/shared"
)

type ToolPlugin interface {
	base.BasePlugin

	GetDefinition() (*shared.ToolDefinition, error)

	Handle(ctx *shared.ToolContext) error
}

package tools

import (
	"fmt"

	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/shared"
)

type RpcServer struct {
	*base.RpcServer
	Impl ToolPlugin
}

func (rs *RpcServer) GetDefinition(_ struct{}, resp *shared.ToolDefinition) error {
	r, err := rs.Impl.GetDefinition()
	if err != nil {
		return fmt.Errorf("error performing server call: %w", err)
	}

	*resp = *r
	return nil
}

func (rs *RpcServer) Handle(ctx *shared.ToolContext, resp *error) error {
	*resp = rs.Impl.Handle(ctx)
	return *resp
}

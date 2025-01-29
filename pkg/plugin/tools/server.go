package tools

import (
	"fmt"

	"github.com/mwantia/queueverse/pkg/plugin/base"
)

type RpcServer struct {
	*base.RpcServer
	Impl ToolPlugin
}

func (rs *RpcServer) GetParameters(_ struct{}, resp *ToolParameters) error {
	r, err := rs.Impl.GetParameters()
	if err != nil {
		return fmt.Errorf("error performing server call: %w", err)
	}

	*resp = *r
	return nil
}

func (rs *RpcServer) Handle(ctx *ToolContext, resp *error) error {
	*resp = rs.Impl.Handle(ctx)
	return *resp
}

package tools

import (
	"fmt"
	"net/rpc"

	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/shared"
)

type RpcClient struct {
	*base.RpcClient
	Client *rpc.Client
}

func (rc *RpcClient) GetDefinition() (*shared.ToolDefinition, error) {
	var resp *shared.ToolDefinition
	if err := rc.Client.Call("Plugin.GetDefinition", struct{}{}, &resp); err != nil {
		return resp, fmt.Errorf("error performing client call: %w", err)
	}

	return resp, nil
}

func (rc *RpcClient) Handle(ctx *shared.ToolContext) error {
	var resp error
	if err := rc.Client.Call("Plugin.Handle", struct{}{}, &resp); err != nil {
		return fmt.Errorf("error performing client call: %w", err)
	}

	return resp
}

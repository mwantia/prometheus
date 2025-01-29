package tools

import (
	"fmt"
	"net/rpc"

	"github.com/mwantia/queueverse/pkg/plugin/base"
)

type RpcClient struct {
	*base.RpcClient
	Client *rpc.Client
}

func (rc *RpcClient) GetParameters() (*ToolParameters, error) {
	var resp *ToolParameters
	if err := rc.Client.Call("Plugin.GetParameters", struct{}{}, &resp); err != nil {
		return resp, fmt.Errorf("error performing client call: %w", err)
	}

	return resp, nil
}

func (rc *RpcClient) Handle(ctx *ToolContext) error {
	var resp error
	if err := rc.Client.Call("Plugin.Handle", struct{}{}, &resp); err != nil {
		return fmt.Errorf("error performing client call: %w", err)
	}

	return resp
}

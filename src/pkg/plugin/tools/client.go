package tools

import (
	"fmt"
	"net/rpc"
)

type RpcClient struct {
	client *rpc.Client
}

func (rc *RpcClient) GetName() (string, error) {
	var resp string
	if err := rc.client.Call("Plugin.GetName", struct{}{}, &resp); err != nil {
		return "", fmt.Errorf("error performing client call: %w", err)
	}

	return resp, nil
}

func (rc *RpcClient) GetParameters() (*ToolParameters, error) {
	var resp *ToolParameters
	if err := rc.client.Call("Plugin.GetParameters", struct{}{}, &resp); err != nil {
		return resp, fmt.Errorf("error performing client call: %w", err)
	}

	return resp, nil
}

func (rc *RpcClient) Handle(ctx *ToolContext) error {
	var resp error
	if err := rc.client.Call("Plugin.Handle", struct{}{}, &resp); err != nil {
		return fmt.Errorf("error performing client call: %w", err)
	}

	return resp
}

func (rc *RpcClient) Probe() error {
	var resp error
	if err := rc.client.Call("Plugin.Probe", struct{}{}, &resp); err != nil {
		return fmt.Errorf("error performing client call: %w", err)
	}

	return resp
}

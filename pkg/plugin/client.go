package plugin

import (
	"fmt"
	"net/rpc"

	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/pkg/plugin/tools"
)

type RpcClient struct {
	client *rpc.Client
}

func (c *RpcClient) GetProvider() (provider.Provider, error) {
	var resp provider.ProviderRPC
	if err := c.client.Call("Plugin.GetProvider", struct{}{}, &resp); err != nil {
		return nil, fmt.Errorf("failed to call rpc: %w", err)
	}

	return &provider.ProviderRPCClient{Client: c.client}, nil
}

func (c *RpcClient) GetTools() (*[]tools.ToolParameters, error) {
	var resp []tools.ToolParameters
	if err := c.client.Call("Plugin.GetTools", struct{}{}, &resp); err != nil {
		return &resp, fmt.Errorf("failed to get plugin info: %w", err)
	}

	return &resp, nil
}

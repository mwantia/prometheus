package identity

import (
	"fmt"
	"net/rpc"
)

type RpcClient struct {
	client *rpc.Client
}

func (c *RpcClient) GetPluginInfo() (*PluginInfo, error) {
	var resp PluginInfo
	if err := c.client.Call("Plugin.GetPluginInfo", struct{}{}, &resp); err != nil {
		return &resp, fmt.Errorf("failed to get plugin info: %w", err)
	}

	return &resp, nil
}

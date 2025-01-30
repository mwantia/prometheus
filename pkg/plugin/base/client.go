package base

import (
	"fmt"
	"net/rpc"
)

type RpcClient struct {
	Client *rpc.Client
}

func (rc *RpcClient) GetPluginInfo() (*PluginInfo, error) {
	var reply *PluginInfo
	if err := rc.Client.Call("Plugin.GetPluginInfo", struct{}{}, &reply); err != nil {
		return nil, fmt.Errorf("error performing client call: %w", err)
	}

	return reply, nil
}

func (rc *RpcClient) GetCapabilities() (*PluginCapabilities, error) {
	var reply *PluginCapabilities
	if err := rc.Client.Call("Plugin.GetCapabilities", struct{}{}, &reply); err != nil {
		return nil, fmt.Errorf("error performing client call: %w", err)
	}

	return reply, nil
}

func (rc *RpcClient) SetConfig(cfg *PluginConfig) error {
	var reply error
	if err := rc.Client.Call("Plugin.SetConfig", cfg, &reply); err != nil {
		return fmt.Errorf("error performing client call: %w", err)
	}

	return reply
}

func (rc *RpcClient) ProbePlugin() error {
	var reply error
	if err := rc.Client.Call("Plugin.ProbePlugin", struct{}{}, &reply); err != nil {
		return fmt.Errorf("error performing client call: %w", err)
	}

	return reply
}

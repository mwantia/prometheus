package plugin

import (
	"errors"
	"fmt"
	"net/rpc"
)

type RpcClient struct {
	Client *rpc.Client
}

func (rc *RpcClient) Name() (string, error) {
	var resp string
	if err := rc.Client.Call("Plugin.Name", struct{}{}, &resp); err != nil {
		return "", fmt.Errorf("failed to get plugin name: %v", err)
	}

	return resp, nil
}

func (rc *RpcClient) GetCapabilities() (PluginCapabilities, error) {
	var resp PluginCapabilities
	if err := rc.Client.Call("Plugin.GetCapabilities", struct{}{}, &resp); err != nil {
		return resp, fmt.Errorf("failed to get plugin capabilities: %v", err)
	}

	return resp, nil
}

func (rc *RpcClient) Setup(s PluginSetup) error {
	var resp error
	if err := rc.Client.Call("Plugin.Setup", s, &resp); err != nil {
		err = errors.Join(err, resp)
		return fmt.Errorf("failed to complete plugin setup: %v", err)
	}

	return resp
}

func (rc *RpcClient) Health() error {
	var resp error
	if err := rc.Client.Call("Plugin.Health", struct{}{}, &resp); err != nil {
		return fmt.Errorf("failed to get plugin health: %v", errors.Join(err, resp))
	}

	return resp
}

func (rc *RpcClient) Cleanup() error {
	var resp error
	if err := rc.Client.Call("Plugin.Cleanup", struct{}{}, &resp); err != nil {
		return fmt.Errorf("failed to complete plugin cleanup: %v", errors.Join(err, resp))
	}

	return resp
}

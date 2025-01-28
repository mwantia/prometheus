package provider

import (
	"net/rpc"
)

type RpcClient struct {
	client *rpc.Client
}

func (rc *RpcClient) Config(cfg map[string]any) error {
	var reply error
	if err := rc.client.Call("Plugin.Config", cfg, &reply); err != nil {
		return err
	}

	return reply
}

func (rc *RpcClient) Chat(req ProviderChatRequest) (ProviderChatResponse, error) {
	var reply ProviderChatResponse
	if err := rc.client.Call("Plugin.Chat", req, &reply); err != nil {
		return reply, err
	}

	return reply, nil
}

func (rc *RpcClient) Probe() error {
	var reply error
	if err := rc.client.Call("Plugin.Probe", struct{}{}, &reply); err != nil {
		return err
	}

	return reply
}

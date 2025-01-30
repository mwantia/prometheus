package provider

import (
	"net/rpc"

	"github.com/mwantia/queueverse/pkg/plugin/base"
)

type RpcClient struct {
	base.RpcClient
	Client *rpc.Client
}

func (rc *RpcClient) GetModels() (*[]ProviderModel, error) {
	var reply *[]ProviderModel
	if err := rc.Client.Call("Plugin.GetModels", struct{}{}, &reply); err != nil {
		return reply, err
	}

	return reply, nil
}

func (rc *RpcClient) Chat(req ProviderChatRequest) (*ProviderChatResponse, error) {
	var reply *ProviderChatResponse
	if err := rc.Client.Call("Plugin.Chat", req, &reply); err != nil {
		return reply, err
	}

	return reply, nil
}

package provider

import (
	"net/rpc"

	"github.com/mwantia/queueverse/pkg/plugin/base"
)

type RpcClient struct {
	base.RpcClient
	Client *rpc.Client
}

func (rc *RpcClient) GetModels() (*[]Model, error) {
	var reply *[]Model
	if err := rc.Client.Call("Plugin.GetModels", struct{}{}, &reply); err != nil {
		return reply, err
	}

	return reply, nil
}

func (rc *RpcClient) Chat(req ChatRequest) (*ChatResponse, error) {
	var reply *ChatResponse
	if err := rc.Client.Call("Plugin.Chat", req, &reply); err != nil {
		return reply, err
	}

	return reply, nil
}

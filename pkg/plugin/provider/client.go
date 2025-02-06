package provider

import (
	"net/rpc"

	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/shared"
)

type RpcClient struct {
	base.RpcClient
	Client *rpc.Client
}

func (rc *RpcClient) GetModels() (*[]shared.Model, error) {
	var reply *[]shared.Model
	if err := rc.Client.Call("Plugin.GetModels", struct{}{}, &reply); err != nil {
		return reply, err
	}

	return reply, nil
}

func (rc *RpcClient) Chat(request shared.ChatRequest, handler shared.ProviderToolHandler) (*shared.ChatResponse, error) {
	args := &ChatArgs{
		Request: request,
		Handler: handler,
	}

	var reply *shared.ChatResponse
	if err := rc.Client.Call("Plugin.Chat", args, &reply); err != nil {
		return reply, err
	}

	return reply, nil
}

func (rc *RpcClient) Embed(req shared.EmbedRequest) (*shared.EmbedResponse, error) {
	var reply *shared.EmbedResponse
	if err := rc.Client.Call("Plugin.Embed", req, &reply); err != nil {
		return reply, err
	}

	return reply, nil
}

package provider

import (
	"context"
	"fmt"
	"net/rpc"

	"github.com/mwantia/queueverse/pkg/plugin/base"
)

type Provider interface {
	base.Base

	Chat(ProviderChatRequest) (ProviderChatResponse, error)
}

type ProviderRPC struct {
	Impl Provider
}

func (p *ProviderRPC) Config(args map[string]any, reply *error) error {
	*reply = p.Impl.Config(args)
	return nil
}

type ChatArgs struct {
	Context context.Context
	Request ProviderChatRequest
}

func (p *ProviderRPC) Chat(args *ChatArgs, reply *ProviderChatResponse) error {
	resp, err := p.Impl.Chat(args.Request)
	if err != nil {
		return err
	}
	*reply = resp
	return nil
}

func (p *ProviderRPC) Probe(args struct{}, reply *error) error {
	*reply = p.Impl.Probe()
	return nil
}

// And your RPC client needs to implement the Provider interface
type ProviderRPCClient struct {
	Client *rpc.Client
}

func (c *RpcClient) GetProvider() (Provider, error) {
	var resp ProviderRPC
	if err := c.client.Call("Plugin.GetProvider", struct{}{}, &resp); err != nil {
		return nil, fmt.Errorf("failed to call rpc: %w", err)
	}

	return &ProviderRPCClient{Client: c.client}, nil
}

// Implement Provider interface on the client side
func (c *ProviderRPCClient) Config(config map[string]any) error {
	var reply error
	err := c.Client.Call("Plugin.Config", config, &reply)
	if err != nil {
		return err
	}
	return reply
}

func (c *ProviderRPCClient) Chat(req ProviderChatRequest) (ProviderChatResponse, error) {
	args := &ChatArgs{
		Request: req,
	}
	var reply ProviderChatResponse
	err := c.Client.Call("Plugin.Chat", args, &reply)
	return reply, err
}

func (c *ProviderRPCClient) Probe() error {
	var reply error
	err := c.Client.Call("Plugin.Probe", struct{}{}, &reply)
	if err != nil {
		return err
	}
	return reply
}

package tools

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type ToolPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Service ToolService
}

func (p *ToolPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RpcServer{
		impl: p.Service,
	}, nil
}

func (p *ToolPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RpcClient{
		client: c,
	}, nil
}

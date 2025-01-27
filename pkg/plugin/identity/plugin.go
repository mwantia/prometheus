package identity

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type IdentityPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Service IdentityService
}

func (p *IdentityPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RpcServer{
		impl: p.Service,
	}, nil
}

func (p *IdentityPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RpcClient{
		client: c,
	}, nil
}

package cache

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type CachePlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl CacheService
}

func (p *CachePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RpcServer{
		impl: p.Impl,
	}, nil
}

func (p *CachePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RpcClient{
		client: c,
	}, nil
}

package plugin

import (
	"net/rpc"

	goplugin "github.com/hashicorp/go-plugin"
)

type PluginDriver struct {
	Impl Plugin
}

func (p *PluginDriver) Server(*goplugin.MuxBroker) (interface{}, error) {
	return &RpcServer{
		Impl: p.Impl,
	}, nil
}

func (p *PluginDriver) Client(b *goplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RpcClient{
		Client: c,
	}, nil
}

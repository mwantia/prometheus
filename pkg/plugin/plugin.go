package plugin

import (
	"net/rpc"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/pkg/plugin/tools"
)

var Handshake = goplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "QUEUEVERSE_PLUGIN",
	MagicCookieValue: "queueverse",
}

var PluginMap = map[string]goplugin.Plugin{
	"plugin": &PluginImpl{},
}

type Plugin interface {
	GetProvider() (provider.Provider, error)

	GetTools() (*[]tools.ToolParameters, error)
}

type PluginImpl struct {
	goplugin.NetRPCUnsupportedPlugin
	Impl Plugin
}

func (p *PluginImpl) Server(*goplugin.MuxBroker) (interface{}, error) {
	return &RpcServer{
		impl: p.Impl,
	}, nil
}

func (p *PluginImpl) Client(b *goplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RpcClient{
		client: c,
	}, nil
}

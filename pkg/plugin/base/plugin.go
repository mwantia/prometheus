package base

import (
	"net/rpc"

	goplugin "github.com/hashicorp/go-plugin"
)

const (
	PluginBaseType     = "base"
	PluginProviderType = "provider"
	PluginToolsType    = "tools"
)

var Handshake = goplugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "QUEUEVERSE_PLUGIN",
	MagicCookieValue: "queueverse",
}

type BasePluginImpl struct {
	goplugin.NetRPCUnsupportedPlugin
	Impl BasePlugin
}

func (impl *BasePluginImpl) Server(*goplugin.MuxBroker) (interface{}, error) {
	return &RpcServer{
		Impl: impl.Impl,
	}, nil
}

func (impl *BasePluginImpl) Client(b *goplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RpcClient{
		Client: c,
	}, nil
}

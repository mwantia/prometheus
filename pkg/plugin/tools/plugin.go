package tools

import (
	"net/rpc"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/mwantia/queueverse/pkg/plugin/base"
)

type ToolPluginImpl struct {
	goplugin.NetRPCUnsupportedPlugin
	Impl ToolPlugin
}

func Serve(plugin ToolPlugin) error {
	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: base.Handshake,
		Plugins: map[string]goplugin.Plugin{
			base.PluginBaseType: &base.BasePluginImpl{
				Impl: plugin,
			},
			base.PluginProviderType: &ToolPluginImpl{
				Impl: plugin,
			},
		},
		GRPCServer: goplugin.DefaultGRPCServer,
	})
	return nil
}

func (impl *ToolPluginImpl) Server(*goplugin.MuxBroker) (interface{}, error) {
	return &RpcServer{
		Impl: impl.Impl,
		RpcServer: &base.RpcServer{
			Impl: impl.Impl,
		},
	}, nil
}

func (*ToolPluginImpl) Client(b *goplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RpcClient{
		Client: c,
		RpcClient: &base.RpcClient{
			Client: c,
		},
	}, nil
}

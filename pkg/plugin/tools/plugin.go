package tools

import (
	"net/rpc"

	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/mwantia/queueverse/pkg/plugin/base"
)

type ToolPluginImpl struct {
	goplugin.NetRPCUnsupportedPlugin
	Impl ToolPlugin
}

func Serve(impl ToolPlugin, logger hclog.Logger) {
	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: base.Handshake,
		Plugins: map[string]goplugin.Plugin{
			base.PluginBaseType: &base.BasePluginImpl{
				Impl: impl,
			},
			base.PluginProviderType: &ToolPluginImpl{
				Impl: impl,
			},
		},
		GRPCServer: goplugin.DefaultGRPCServer,
		Logger:     logger,
	})
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

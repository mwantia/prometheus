package provider

import (
	"net/rpc"

	hclog "github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/mwantia/queueverse/pkg/plugin/base"
)

type ProviderPluginImpl struct {
	goplugin.NetRPCUnsupportedPlugin
	Impl ProviderPlugin
}

func Serve(impl ProviderPlugin, logger hclog.Logger) error {
	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: base.Handshake,
		Plugins: map[string]goplugin.Plugin{
			base.PluginBaseType: &base.BasePluginImpl{
				Impl: impl,
			},
			base.PluginProviderType: &ProviderPluginImpl{
				Impl: impl,
			},
		},
		GRPCServer: goplugin.DefaultGRPCServer,
		Logger:     logger,
	})
	return nil
}

func (impl *ProviderPluginImpl) Server(*goplugin.MuxBroker) (interface{}, error) {
	return &RpcServer{
		Impl: impl.Impl,
		RpcServer: base.RpcServer{
			Impl: impl.Impl,
		},
	}, nil
}

func (*ProviderPluginImpl) Client(b *goplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RpcClient{
		Client: c,
		RpcClient: base.RpcClient{
			Client: c,
		},
	}, nil
}

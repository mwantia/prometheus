package provider

import (
	"github.com/hashicorp/go-plugin"
	"github.com/mwantia/queueverse/pkg/plugin/base"
)

type ProviderImpl struct {
	plugin.NetRPCUnsupportedPlugin
	Impl Provider
}

func Serve(impl Provider) error {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: base.Handshake,
		Plugins: map[string]plugin.Plugin{
			base.PluginBaseType: &base.BaseImpl{
				Impl: impl,
			},
			base.PluginProviderType: &ProviderImpl{
				Impl: impl,
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
	return nil
}

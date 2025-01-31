package mock

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
)

const (
	PluginType    = base.PluginProviderType
	PluginName    = "mock"
	PluginAuthor  = "mwantia"
	PluginVersion = "v0.0.1"
)

type MockProvider struct {
	provider.UnimplementedProviderPlugin

	Context context.Context
	Logger  hclog.Logger
}

func (*MockProvider) GetPluginInfo() (*base.PluginInfo, error) {
	return &base.PluginInfo{
		Type:    PluginType,
		Name:    PluginName,
		Author:  PluginAuthor,
		Version: PluginVersion,
	}, nil
}

func (*MockProvider) SetConfig(*base.PluginConfig) error {
	return nil
}

func (*MockProvider) ProbePlugin() error {
	return nil
}

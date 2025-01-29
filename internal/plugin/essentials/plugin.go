package essentials

import (
	"context"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/mitchellh/mapstructure"
	"github.com/mwantia/queueverse/pkg/plugin"
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
)

const (
	PluginType    = base.PluginProviderType
	PluginName    = "essentials"
	PluginVersion = "v0.0.1-alpha"
)

type EssentialsProviderPlugin struct {
	provider.UnimplementedProviderPlugin

	Cfg    *EssentialsProviderConfig
	Logger hclog.Logger
}

type EssentialsProviderConfig struct {
	Foo string `mapstructure:"foo,omitempty"`
}

func Serve() error {
	return plugin.ServeContext(func(ctx context.Context, logger hclog.Logger) interface{} {
		return &EssentialsProviderPlugin{
			Logger: logger,
		}
	})
}

func (*EssentialsProviderPlugin) GetPluginInfo() (*base.PluginInfo, error) {
	return &base.PluginInfo{
		Type:    PluginType,
		Name:    PluginName,
		Version: PluginVersion,
	}, nil
}

func (p *EssentialsProviderPlugin) SetConfig(cfg *base.PluginConfig) error {
	if err := mapstructure.Decode(cfg.ConfigMap, &p.Cfg); err != nil {
		return err
	}

	p.Logger.Info("Foo: " + p.Cfg.Foo)

	return nil
}

func (*EssentialsProviderPlugin) ProbePlugin() error {
	return nil
}

func (*EssentialsProviderPlugin) Chat(req provider.ProviderChatRequest) (*provider.ProviderChatResponse, error) {
	return &provider.ProviderChatResponse{
		Model:   req.Model,
		Message: req.Messages[0],
	}, nil
}

package anthropic

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/liushuangls/go-anthropic/v2"
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
)

const (
	PluginType    = base.PluginProviderType
	PluginName    = "anthropic"
	PluginAuthor  = "mwantia"
	PluginVersion = "v0.0.1"
)

type AnthropicProvider struct {
	provider.UnimplementedProviderPlugin

	Context context.Context
	Logger  hclog.Logger
	Config  AnthropicConfig
	Client  *anthropic.Client
}

func (*AnthropicProvider) GetPluginInfo() (*base.PluginInfo, error) {
	return &base.PluginInfo{
		Type:    PluginType,
		Name:    PluginName,
		Author:  PluginAuthor,
		Version: PluginVersion,
	}, nil
}

func (p *AnthropicProvider) SetConfig(cfg *base.PluginConfig) error {
	if err := p.setConfig(cfg.ConfigMap); err != nil {
		return err
	}

	if p.Config.Token == "" {
		return fmt.Errorf("config 'token' must be set")
	}

	p.Client = anthropic.NewClient(p.Config.Token)

	return nil
}

func (p *AnthropicProvider) ProbePlugin() error {
	return nil
}

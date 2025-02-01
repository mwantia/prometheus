package ollama

import (
	"context"
	"fmt"
	"net/http"

	hclog "github.com/hashicorp/go-hclog"

	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/plugins/ollama/api"
)

const (
	PluginType    = base.PluginProviderType
	PluginName    = "ollama"
	PluginAuthor  = "mwantia"
	PluginVersion = "v0.0.1"
)

type OllamaProvider struct {
	provider.UnimplementedProviderPlugin

	Context context.Context
	Logger  hclog.Logger
	Config  OllamaProviderConfig
	Client  api.Client
}

func (*OllamaProvider) GetCapabilities() (*base.PluginCapabilities, error) {
	return &base.PluginCapabilities{
		Types: []base.PluginCapabilityType{
			base.Generate,
			base.Embed,
		},
	}, nil
}

func (*OllamaProvider) GetPluginInfo() (*base.PluginInfo, error) {
	return &base.PluginInfo{
		Type:    PluginType,
		Name:    PluginName,
		Author:  PluginAuthor,
		Version: PluginVersion,
	}, nil
}

func (p *OllamaProvider) SetConfig(cfg *base.PluginConfig) error {
	if err := p.setConfig(cfg.ConfigMap); err != nil {
		return err
	}

	if p.Config.Endpoint == "" {
		return fmt.Errorf("config 'endpoint' is not defined")
	}

	p.Client = api.CreateClient(http.DefaultClient, api.ClientConfig{
		Endpoint: p.Config.Endpoint,
	})

	return nil
}

func (p *OllamaProvider) ProbePlugin() error {
	if p.Client == nil {
		return fmt.Errorf("ollama client is undefined")
	}

	ok, err := p.Client.Health(p.Context)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("client status is 'unhealthy'")
	}

	return nil
}

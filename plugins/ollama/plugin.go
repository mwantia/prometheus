package ollama

import (
	"context"
	"fmt"
	"net/http"

	hclog "github.com/hashicorp/go-hclog"

	"github.com/mwantia/queueverse/pkg/plugin"
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/plugins/ollama/api"
)

const (
	PluginType    = base.PluginProviderType
	PluginName    = "ollama"
	PluginVersion = "v0.0.1"
)

type OllamaPlugin struct {
	provider.UnimplementedProviderPlugin
	Context context.Context
	Logger  hclog.Logger
	Config  OllamaPluginConfig
	Client  api.Client
}

func Serve() error {
	return plugin.ServeContext(func(ctx context.Context, logger hclog.Logger) interface{} {
		return &OllamaPlugin{
			Context: ctx,
			Logger:  logger,
		}
	})
}

func (*OllamaPlugin) GetPluginInfo() (*base.PluginInfo, error) {
	return &base.PluginInfo{
		Type:    PluginType,
		Name:    PluginName,
		Version: PluginVersion,
	}, nil
}

func (p *OllamaPlugin) SetConfig(cfg *base.PluginConfig) error {
	if err := p.setConfig(cfg.ConfigMap); err != nil {
		return err
	}

	if p.Config.Endpoint == "" {
		return fmt.Errorf("config 'endpoint' is not defined")
	}

	p.Client = api.CreateClient(http.DefaultClient, api.ClientConfig{
		Endpoint: p.Config.Endpoint,
		Model:    p.Config.Model,
	})

	return nil
}

func (p *OllamaPlugin) ProbePlugin() error {
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

package essentials

import (
	"context"

	"github.com/mwantia/queueverse/pkg/plugin"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/pkg/plugin/tools"
)

type EssentialsService struct {
	tools *[]tools.ToolParameters
}

func Serve() error {
	return plugin.ServeContextMux(plugin.PluginContextFactoryMuxMap{
		"provider": func(ctx context.Context) interface{} {
			if err := ctx.Err(); err != nil {
				return nil
			}

			return nil
		},
	})
}

func (s *EssentialsService) GetTools() (*[]tools.ToolParameters, error) {
	return s.tools, nil
}

func (s *EssentialsService) GetProvider() (provider.Provider, error) {
	return &EssentialsProvider{}, nil
}

type EssentialsProvider struct{}

func (p *EssentialsProvider) Config(map[string]any) error {
	return nil
}

func (p *EssentialsProvider) Chat(req provider.ProviderChatRequest) (provider.ProviderChatResponse, error) {
	return provider.ProviderChatResponse{
		Model:    req.Model,
		Message:  req.Messages[0],
		Metadata: map[string]any{},
	}, nil
}

func (p *EssentialsProvider) Probe() error {
	return nil
}

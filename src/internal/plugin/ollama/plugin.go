package ollama

import (
	"context"
	"fmt"
	"log"

	"github.com/mwantia/prometheus/pkg/msg"
	"github.com/mwantia/prometheus/pkg/plugin"
)

type OllamaPlugin struct {
	plugin.DefaultPlugin

	Context context.Context
	Config  OllamaConfig
	Hub     msg.MessageHub
}

func NewPlugin() *OllamaPlugin {
	return &OllamaPlugin{
		Context: context.Background(),
	}
}

func (p *OllamaPlugin) Name() (string, error) {
	return "ollama", nil
}

func (p *OllamaPlugin) Setup(s plugin.PluginSetup) error {
	p.Hub = msg.NewMessageHubCacher(s.Hub)

	if err := p.loadConfig(s.Data); err != nil {
		log.Printf("Error converting mapstructure: %v", err)
	}

	return nil
}

func (p *OllamaPlugin) Health() error {
	return nil
}

func (p *OllamaPlugin) Cleanup() error {
	if err := p.Hub.Cleanup(); err != nil {
		return fmt.Errorf("unable to cleanup message hub manager: %v", err)
	}

	return nil
}

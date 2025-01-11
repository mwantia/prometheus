package debug

import (
	"context"
	"log"

	"github.com/mwantia/prometheus/pkg/plugin"
)

type DebugPlugin struct {
	plugin.DefaultPlugin

	Context context.Context
	Config  *DebugConfig
}

func NewPlugin() *DebugPlugin {
	return &DebugPlugin{
		Context: context.Background(),
	}
}

func (p *DebugPlugin) Name() (string, error) {
	return "debug", nil
}

func (p *DebugPlugin) Setup(s plugin.PluginSetup) error {
	if err := p.loadConfig(s.Data); err != nil {
		log.Printf("Error converting mapstructure: %v", err)
	}

	log.Printf("Foo: %s", p.Config.Foo)

	return nil
}

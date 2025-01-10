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

	producer, err := s.Hub.CreateProducer("global")
	if err != nil {
		log.Printf("Error completing setup for message hub: %v", err)
	}

	if err := producer.Write(p.Context, "Hello World"); err != nil {
		log.Printf("Error writing message: %v", err)
	}
	log.Printf("Foo: %s", p.Config.Foo)

	return nil
}

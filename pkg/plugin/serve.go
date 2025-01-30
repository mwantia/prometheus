package plugin

import (
	"context"
	"fmt"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/pkg/plugin/tools"
)

type PluginFactory func(log hclog.Logger) interface{}

type PluginContextFactory func(ctx context.Context, log hclog.Logger) interface{}

func Serve(pf PluginFactory) error {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		JSONFormat: true,
	})

	plugin := pf(logger)
	return servePlugin(plugin, logger)
}

func ServeContext(pcf PluginContextFactory) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		JSONFormat: true,
	})

	plugin := pcf(ctx, logger)
	return servePlugin(plugin, logger)
}

func servePlugin(plugin interface{}, logger hclog.Logger) error {
	return func() error {
		switch impl := plugin.(type) {
		case provider.ProviderPlugin:
			return provider.Serve(impl, logger)
		case tools.ToolPlugin:
			return nil
		default:
			return fmt.Errorf(`unsupported plugin type.
			Ensure that all funcs are defined correctly, even in defaults.go.
			`)
		}
	}()
}

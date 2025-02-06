package plugin

import (
	"context"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/pkg/plugin/tools"
)

type PluginFactory func(log hclog.Logger) interface{}

type PluginContextFactory func(ctx context.Context, log hclog.Logger) interface{}

func Serve(pf PluginFactory) {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		JSONFormat: true,
	})

	plugin := pf(logger)
	servePlugin(plugin, logger)
}

func ServeContext(pcf PluginContextFactory) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		JSONFormat: true,
	})

	plugin := pcf(ctx, logger)
	servePlugin(plugin, logger)
}

func servePlugin(plugin interface{}, logger hclog.Logger) {
	switch impl := plugin.(type) {
	case provider.ProviderPlugin:
		provider.Serve(impl, logger)
	case tools.ToolPlugin:
		tools.Serve(impl, logger)
	default:
		panic(`unsupported plugin type. Ensure that all funcs are defined correctly, even in defaults.go.`)
	}
}

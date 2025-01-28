package plugin

import (
	"context"
	"fmt"
	"os"

	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/pkg/plugin/tools"
)

type PluginFactory func() interface{}

type PluginFactoryMuxMap map[string]PluginFactory

type PluginContextFactory func(ctx context.Context) interface{}

type PluginContextFactoryMuxMap map[string]PluginContextFactory

func Serve(pf PluginFactory) error {
	plugin := pf()
	return servePlugin(plugin)
}

func ServeMux(mux PluginFactoryMuxMap) error {
	if len(os.Args) != 2 {
		return fmt.Errorf("only one additional argument expected for 'os.Args'")
	}
	pf, ok := mux[os.Args[1]]
	if !ok {
		return fmt.Errorf("unknown plugin: %s", os.Args[1])
	}

	return Serve(pf)
}

func ServeContext(pcf PluginContextFactory) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	plugin := pcf(ctx)
	return servePlugin(plugin)
}

func ServeContextMux(mux PluginContextFactoryMuxMap) error {
	if len(os.Args) != 2 {
		return fmt.Errorf("only one additional argument expected for 'os.Args'")
	}
	pf, ok := mux[os.Args[1]]
	if !ok {
		return fmt.Errorf("unknown plugin: %s", os.Args[1])
	}

	return ServeContext(pf)
}

func servePlugin(plugin interface{}) error {
	return func() error {
		switch impl := plugin.(type) {
		case provider.Provider:
			return provider.Serve(impl)
		case tools.ToolService:
			return nil
		default:
			return fmt.Errorf("unsupported plugin type")
		}
	}()
}

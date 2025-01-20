package debug

import (
	"github.com/mwantia/prometheus/pkg/plugin"
	"github.com/mwantia/prometheus/pkg/plugin/identity"
	"github.com/mwantia/prometheus/pkg/plugin/tools"
)

type DebugIdentityService struct {
	identity.DefaultUnimplementedService
}

type GetDebugTool struct {
	tools.DefaultUnimplementedService
}

const DebugToolName = "debug"

func Serve() error {
	return plugin.ServeTools(&DebugIdentityService{}, []tools.ToolService{
		&GetDebugTool{},
	})
}

func (p *DebugIdentityService) GetPluginInfo() (*identity.PluginInfo, error) {
	return &identity.PluginInfo{
		Name:    DebugToolName,
		Version: "0.0.1",
		Author:  "mwantia",
		Services: []identity.PluginServiceInfo{
			{
				Name:        DebugToolName,
				Type:        identity.ToolServiceType,
				Description: "Tool used to receive and display the current debug information",
			},
		},
	}, nil
}

func (t *GetDebugTool) GetName() (string, error) {
	return DebugToolName, nil
}

func (t *GetDebugTool) GetParameters() (*tools.ToolParameters, error) {
	return &tools.ToolParameters{
		ReturnType:  "string",
		Description: "Returns a list of all available debug informations.",
	}, nil
}

func (t *GetDebugTool) Handle(ctx *tools.ToolContext) error {
	return nil
}

func (t *GetDebugTool) Probe() error {
	return nil
}

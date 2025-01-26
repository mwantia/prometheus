package essentials

import (
	"github.com/mwantia/queueverse/pkg/plugin"
	"github.com/mwantia/queueverse/pkg/plugin/identity"
	"github.com/mwantia/queueverse/pkg/plugin/tools"
)

type EssentialsIdentityService struct {
	identity.DefaultUnimplementedService
}

func Serve() error {
	return plugin.ServeTools(&EssentialsIdentityService{}, []tools.ToolService{
		&GetCurrentTimeTool{},
	})
}

func (srv *EssentialsIdentityService) GetPluginInfo() (*identity.PluginInfo, error) {
	return &identity.PluginInfo{
		Name:    "essentials",
		Version: "0.0.1",
		Author:  "mwantia",
		Services: []identity.PluginServiceInfo{
			{
				Name:        "get_current_time",
				Type:        identity.ToolServiceType,
				Description: "Tool used to get the current time in the specified timezone",
			},
		},
	}, nil
}

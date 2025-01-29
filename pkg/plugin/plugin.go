package plugin

import (
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/pkg/plugin/tools"
)

var Plugins = map[string]goplugin.Plugin{
	base.PluginBaseType:     &base.BasePluginImpl{},
	base.PluginProviderType: &provider.ProviderPluginImpl{},
	base.PluginToolsType:    &tools.ToolPluginImpl{},
}

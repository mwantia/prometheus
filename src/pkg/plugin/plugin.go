package plugin

import (
	"github.com/hashicorp/go-plugin"
	"github.com/mwantia/prometheus/pkg/plugin/cache"
	"github.com/mwantia/prometheus/pkg/plugin/identity"
	"github.com/mwantia/prometheus/pkg/plugin/tools"
)

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "PROMETHEUS_PLUGIN",
	MagicCookieValue: "prometheus",
}

var PluginMap = map[string]plugin.Plugin{
	"identity": &identity.IdentityPlugin{},
	"tool.0":   &tools.ToolPlugin{},
	"tool.1":   &tools.ToolPlugin{},
	"tool.2":   &tools.ToolPlugin{},
	"tool.3":   &tools.ToolPlugin{},
	"tool.4":   &tools.ToolPlugin{},
	"cache":    &cache.CachePlugin{},
}

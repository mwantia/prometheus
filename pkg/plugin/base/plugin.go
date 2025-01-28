package base

import "github.com/hashicorp/go-plugin"

const (
	PluginBaseType     = "base"
	PluginProviderType = "provider"
	PluginToolsType    = "tools"
)

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "QUEUEVERSE_PLUGIN",
	MagicCookieValue: "queueverse",
}

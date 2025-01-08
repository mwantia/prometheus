package plugin

import (
	"github.com/hashicorp/go-plugin"
	"github.com/mwantia/prometheus/pkg/msg"
)

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "PROMETHEUS",
	MagicCookieValue: "prometheus",
}

var Plugins = map[string]plugin.Plugin{
	"driver": &PluginDriver{},
}

type Plugin interface {
	Name() (string, error)

	GetCapabilities() (PluginCapabilities, error)

	Setup(s PluginSetup) error

	Health() error

	Cleanup() error
}

type PluginSetup struct {
	Hub  msg.MessageHub
	Data map[string]interface{}
}

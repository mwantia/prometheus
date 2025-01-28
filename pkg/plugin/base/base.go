package base

import "github.com/hashicorp/go-plugin"

type Base interface {
	Config(map[string]any) error

	Probe() error
}

type BaseImpl struct {
	plugin.NetRPCUnsupportedPlugin
	Impl Base
}

package plugin

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DefaultPlugin struct{}

func (p *DefaultPlugin) Name() (string, error) {
	return "", status.Error(codes.Unimplemented, "Name() not implemented")
}

func (p *DefaultPlugin) GetCapabilities() (PluginCapabilities, error) {
	return PluginCapabilities{}, nil
}

func (p *DefaultPlugin) Setup(d map[string]interface{}) error {
	return status.Error(codes.Unimplemented, "Setup() not implemented")
}

func (p *DefaultPlugin) Health() error {
	return nil
}

func (p *DefaultPlugin) Cleanup() error {
	return nil
}

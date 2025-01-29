package base

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UnimplementedBasePlugin struct{}

func (p *UnimplementedBasePlugin) GetPluginInfo() (*PluginInfo, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

func (p *UnimplementedBasePlugin) SetConfig(*PluginConfig) error {
	return status.Error(codes.Unimplemented, "Not implemented")
}

func (p *UnimplementedBasePlugin) ProbePlugin() error {
	return status.Error(codes.Unimplemented, "Not implemented")
}

package base

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UnimplementedBasePlugin struct{}

func (*UnimplementedBasePlugin) GetCapabilities() (*PluginCapabilities, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

func (*UnimplementedBasePlugin) GetPluginInfo() (*PluginInfo, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

func (*UnimplementedBasePlugin) SetConfig(*PluginConfig) error {
	return status.Error(codes.Unimplemented, "Not implemented")
}

func (*UnimplementedBasePlugin) ProbePlugin() error {
	return status.Error(codes.Unimplemented, "Not implemented")
}

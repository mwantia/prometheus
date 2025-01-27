package identity

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DefaultNilService struct{}

func (s *DefaultNilService) GetPluginInfo() (PluginInfo, error) {
	return PluginInfo{}, nil
}

type DefaultUnimplementedService struct{}

func (s *DefaultUnimplementedService) GetPluginInfo() (*PluginInfo, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

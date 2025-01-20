package tools

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DefaultNilService struct{}

func (s *DefaultNilService) GetName() (*string, error) {
	return nil, nil
}

func (s *DefaultNilService) GetParameters() (ToolParameters, error) {
	return ToolParameters{}, nil
}

func (s *DefaultNilService) Handle(ctx ToolContext) error {
	return nil
}

func (d *DefaultNilService) Probe() error {
	return nil
}

type DefaultUnimplementedService struct{}

func (s *DefaultUnimplementedService) GetName() (*string, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

func (s *DefaultUnimplementedService) GetParameters() (*ToolParameters, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

func (s *DefaultUnimplementedService) Handle(ctx ToolContext) error {
	return status.Error(codes.Unimplemented, "Not implemented")
}

func (d *DefaultUnimplementedService) Probe() error {
	return status.Error(codes.Unimplemented, "Not implemented")
}

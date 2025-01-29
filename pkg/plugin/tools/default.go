package tools

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UnimplementedToolPlugin struct{}

func (*UnimplementedToolPlugin) GetName() (*string, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

func (*UnimplementedToolPlugin) GetParameters() (*ToolParameters, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

func (*UnimplementedToolPlugin) Handle(ctx ToolContext) error {
	return status.Error(codes.Unimplemented, "Not implemented")
}

func (*UnimplementedToolPlugin) Probe() error {
	return status.Error(codes.Unimplemented, "Not implemented")
}

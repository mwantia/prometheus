package tools

import (
	"github.com/mwantia/queueverse/pkg/plugin/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UnimplementedToolPlugin struct{}

func (*UnimplementedToolPlugin) GetName() (*string, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

func (*UnimplementedToolPlugin) GetDefinition() (*shared.ToolDefinition, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

func (*UnimplementedToolPlugin) Handle(ctx shared.ToolContext) error {
	return status.Error(codes.Unimplemented, "Not implemented")
}

func (*UnimplementedToolPlugin) Probe() error {
	return status.Error(codes.Unimplemented, "Not implemented")
}

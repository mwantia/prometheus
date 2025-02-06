package shared

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProviderToolHandler interface {
	GetTools() []ToolDefinition

	Execute(context.Context, ToolFunction) (string, error)
}

type UnimplementedProviderToolHandler struct{}

func (*UnimplementedProviderToolHandler) GetTools() []ToolDefinition {
	return make([]ToolDefinition, 0)
}

func (*UnimplementedProviderToolHandler) Execute(context.Context, ToolFunction) (string, error) {
	return "", status.Error(codes.Unimplemented, "Not implemented")
}

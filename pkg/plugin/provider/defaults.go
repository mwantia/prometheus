package provider

import (
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UnimplementedProviderPlugin struct {
	base.UnimplementedBasePlugin
}

func (*UnimplementedProviderPlugin) GetModels() (*[]shared.Model, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

func (*UnimplementedProviderPlugin) Chat(shared.ChatRequest, shared.ProviderToolHandler) (*shared.ChatResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

func (*UnimplementedProviderPlugin) Embed(shared.EmbedRequest) (*shared.EmbedResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

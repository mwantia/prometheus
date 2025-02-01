package provider

import (
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UnimplementedProviderPlugin struct {
	base.UnimplementedBasePlugin
}

func (*UnimplementedProviderPlugin) GetModels() (*[]Model, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

func (*UnimplementedProviderPlugin) Chat(ChatRequest) (*ChatResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

func (*UnimplementedProviderPlugin) Embed(EmbedRequest) (*EmbedResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

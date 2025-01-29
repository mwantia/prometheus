package provider

import (
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UnimplementedProviderPlugin struct {
	base.UnimplementedBasePlugin
}

func (p *UnimplementedProviderPlugin) Chat(ProviderChatRequest) (*ProviderChatResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}

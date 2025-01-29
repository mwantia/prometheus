package provider

import (
	"context"

	"github.com/mwantia/queueverse/pkg/plugin/base"
)

type ProviderPlugin interface {
	base.BasePlugin

	Chat(ProviderChatRequest) (*ProviderChatResponse, error)
}

type ChatArgs struct {
	Context context.Context
	Request ProviderChatRequest
}

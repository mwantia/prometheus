package provider

import (
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/shared"
)

type ProviderPlugin interface {
	base.BasePlugin

	GetModels() (*[]shared.Model, error)

	Chat(shared.ChatRequest, shared.ProviderToolHandler) (*shared.ChatResponse, error)

	Embed(shared.EmbedRequest) (*shared.EmbedResponse, error)
}

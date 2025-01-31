package provider

import (
	"github.com/mwantia/queueverse/pkg/plugin/base"
)

type ProviderPlugin interface {
	base.BasePlugin

	GetModels() (*[]Model, error)

	Chat(ChatRequest) (*ChatResponse, error)
}

package provider

import (
	"github.com/mwantia/queueverse/pkg/plugin/base"
)

type ProviderPlugin interface {
	base.BasePlugin

	GetModels() (*[]ProviderModel, error)

	Chat(ProviderChatRequest) (*ProviderChatResponse, error)
}

type ProviderModel struct {
	Name     string         `json:"name"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type ProviderChatRequest struct {
	Model    string                `json:"model"`
	Messages []ProviderChatMessage `json:"messages"`
	Metadata map[string]any        `json:"metadata,omitempty"`
}

type ProviderChatMessage struct {
	ID      string `json:"id"`
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ProviderChatResponse struct {
	Model    string              `json:"model"`
	Message  ProviderChatMessage `json:"message"`
	Metadata map[string]any      `json:"metadata,omitempty"`
}

package ollama

import (
	"strings"

	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/plugins/ollama/api"
)

func (p *OllamaPlugin) Chat(req provider.ProviderChatRequest) (*provider.ProviderChatResponse, error) {
	var text strings.Builder

	if err := p.Client.Chat(p.Context, api.ChatRequest{
		Model: req.Model,
		Messages: []api.ChatMessage{
			{
				Role:    req.Messages[0].Role,
				Content: req.Messages[0].Content,
			},
		},
		KeepAlive:   -1,
		ContextSize: 8192,
	}, func(resp api.ChatResponse) error {
		_, err := text.WriteString(resp.Message.Content)
		return err
	}); err != nil {
		return nil, err
	}

	return &provider.ProviderChatResponse{
		Model: req.Model,
		Message: provider.ProviderChatMessage{
			Role:    "assistant",
			Content: text.String(),
		},
	}, nil
}

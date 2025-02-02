package ollama

import (
	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/plugins/ollama/api"
)

func (p *OllamaProvider) GetModels() (*[]provider.Model, error) {
	tags, err := p.Client.Tags(p.Context)
	if err != nil {
		return nil, err
	}

	resp := make([]provider.Model, 0)
	for _, tag := range tags {
		resp = append(resp, provider.Model{
			Name: tag.Name,
			Metadata: map[string]any{
				"size":   tag.Size,
				"digest": tag.Digest,
			},
		})
	}
	return &resp, nil
}

func (p *OllamaProvider) Chat(req provider.ChatRequest) (*provider.ChatResponse, error) {
	result := provider.ChatResponse{
		Model:    req.Model,
		Messages: make([]provider.ChatMessage, 0),
	}

	if err := p.Client.Chat(p.Context, CreateMessageRequest(req), func(resp api.ChatResponse) error {
		message := provider.ChatMessage{
			Role:      resp.Message.Role,
			Content:   resp.Message.Content,
			ToolCalls: make([]provider.ToolCall, 0),
		}

		for _, tc := range resp.Message.ToolCalls {
			message.ToolCalls = append(message.ToolCalls, provider.ToolCall{
				Function: provider.ToolFunction{
					Index:     tc.Function.Index,
					Name:      tc.Function.Name,
					Arguments: tc.Function.Arguments,
				},
			})
		}

		result.Messages = append(result.Messages, message)
		return nil
	}); err != nil {
		return nil, err
	}

	return &result, nil
}

func (p *OllamaProvider) Embed(req provider.EmbedRequest) (*provider.EmbedResponse, error) {
	var embeddings [][]float32

	if err := p.Client.Embed(p.Context, api.EmbedRequest{
		Model: req.Model,
		Input: req.Input,
	}, func(resp api.EmbedResponse) error {
		embeddings = resp.Embeddings
		return nil
	}); err != nil {
		return nil, err
	}

	return &provider.EmbedResponse{
		Model:      req.Model,
		Embeddings: embeddings,
	}, nil
}

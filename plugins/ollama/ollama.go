package ollama

import (
	"strings"

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

	return &provider.ChatResponse{
		Model:   req.Model,
		Message: provider.AssistantMessage(text.String()),
	}, nil
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

package anthropic

import (
	"strings"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p *AnthropicProvider) GetModels() (*[]provider.Model, error) {
	return &[]provider.Model{
		{
			Name: string(anthropic.ModelClaude3Dot5HaikuLatest),
			Metadata: map[string]any{
				"size": 0,
			},
		},
		{
			Name: string(anthropic.ModelClaude3Dot5SonnetLatest),
			Metadata: map[string]any{
				"size": 0,
			},
		},
	}, nil
}

func (p *AnthropicProvider) Chat(req provider.ChatRequest) (*provider.ChatResponse, error) {
	resp, err := p.Client.CreateMessages(p.Context, anthropic.MessagesRequest{
		Model: anthropic.Model(req.Model),
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(req.Messages[0].Content),
		},
		MaxTokens: 1000,
	})
	if err != nil {
		return nil, err
	}

	var text strings.Builder
	for _, content := range resp.Content {
		text.WriteString(content.GetText())
	}

	return &provider.ChatResponse{
		Model:   req.Model,
		Message: provider.AssistantMessage(text.String()),
	}, nil
}

func (*AnthropicProvider) Embed(provider.EmbedRequest) (*provider.EmbedResponse, error) {
	return nil, status.Error(codes.Unavailable, "Embed models are not supported")
}

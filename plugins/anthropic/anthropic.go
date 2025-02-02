package anthropic

import (
	"fmt"

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
	resp, err := p.Client.CreateMessages(p.Context, CreateMessageRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to create messages: %w", err)
	}

	return CreateMessageResponse(resp)
}

func (*AnthropicProvider) Embed(provider.EmbedRequest) (*provider.EmbedResponse, error) {
	return nil, status.Error(codes.Unavailable, "Embed models are not supported")
}

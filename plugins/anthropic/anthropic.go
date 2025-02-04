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
			Name:     string(anthropic.ModelClaude3Dot5HaikuLatest),
			Metadata: map[string]any{},
		},
		{
			Name:     string(anthropic.ModelClaude3Dot5SonnetLatest),
			Metadata: map[string]any{},
		},
	}, nil
}

func (p *AnthropicProvider) Chat(input provider.ChatRequest) (*provider.ChatResponse, error) {
	request := anthropic.MessagesRequest{
		Model: anthropic.Model(input.Model),
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(input.Message.Content),
		},
	}

	response, err := p.Client.CreateMessages(p.Context, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create messages: %w", err)
	}

	return &provider.ChatResponse{
		Model: string(response.Model),
		Message: provider.Message{
			Content: response.Content[0].GetText(),
		},
	}, nil
}

func (*AnthropicProvider) Embed(provider.EmbedRequest) (*provider.EmbedResponse, error) {
	return nil, status.Error(codes.Unavailable, "Embed models are not supported")
}

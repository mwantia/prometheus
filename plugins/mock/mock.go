package mock

import (
	"strings"

	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (*MockProvider) GetModels() (*[]provider.Model, error) {
	return &[]provider.Model{
		{
			Name: MockLoremIpsumModel08,
			Metadata: map[string]any{
				"content_length": 8,
			},
		},
		{
			Name: MockLoremIpsumModel16,
			Metadata: map[string]any{
				"content_length": 16,
			},
		},
		{
			Name: MockLoremIpsumModel32,
			Metadata: map[string]any{
				"content_length": 32,
			},
		},
	}, nil
}

func (*MockProvider) Chat(req provider.ChatRequest) (*provider.ChatResponse, error) {
	var text strings.Builder

	switch req.Model {
	case MockLoremIpsumModel08:
		for i := range 8 {
			text.WriteString(MockLoremIpsumContents[i])
		}
	case MockLoremIpsumModel16:
		for i := range 16 {
			text.WriteString(MockLoremIpsumContents[i])
		}
	case MockLoremIpsumModel32:
		for i := range 32 {
			text.WriteString(MockLoremIpsumContents[i])
		}
	}

	return &provider.ChatResponse{
		Model: req.Model,
		Messages: []provider.ChatMessage{
			provider.AssistantMessage(text.String()),
		},
		Metadata: req.Metadata,
	}, nil
}

func (*MockProvider) Embed(provider.EmbedRequest) (*provider.EmbedResponse, error) {
	return nil, status.Error(codes.Unavailable, "Embed models are not supported")
}

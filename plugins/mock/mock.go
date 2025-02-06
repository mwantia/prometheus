package mock

import (
	"strings"

	"github.com/mwantia/queueverse/pkg/plugin/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (*MockProvider) GetModels() (*[]shared.Model, error) {
	return &[]shared.Model{
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

func (p *MockProvider) Chat(req shared.ChatRequest, handler shared.ProviderToolHandler) (*shared.ChatResponse, error) {
	var text strings.Builder

	result, _ := handler.Execute(p.Context, shared.ToolFunction{
		Name: "mock",
		Arguments: map[string]any{
			"foo": "bar",
		},
	})
	p.Logger.Warn(result)

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

	return &shared.ChatResponse{
		Model: req.Model,
		Message: shared.Message{
			Content: text.String(),
		},
		Metadata: req.Metadata,
	}, nil
}

func (*MockProvider) Embed(shared.EmbedRequest) (*shared.EmbedResponse, error) {
	return nil, status.Error(codes.Unavailable, "Embed models are not supported")
}

package mock

import "github.com/mwantia/queueverse/pkg/plugin/provider"

func (*MockProvider) GetModels() (*[]provider.Model, error) {
	return &[]provider.Model{
		{
			Name: MockModelName,
			Metadata: map[string]any{
				"size": 0,
			},
		},
	}, nil
}

func (*MockProvider) Chat(req provider.ChatRequest) (*provider.ChatResponse, error) {
	return &provider.ChatResponse{
		Model:   MockModelName,
		Message: provider.AssistantMessage(MockResponseContent),
	}, nil
}

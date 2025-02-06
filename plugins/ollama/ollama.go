package ollama

import (
	"github.com/mwantia/queueverse/pkg/plugin/shared"
	"github.com/mwantia/queueverse/plugins/ollama/api"
)

func (p *OllamaProvider) GetModels() (*[]shared.Model, error) {
	tags, err := p.Client.Tags(p.Context)
	if err != nil {
		return nil, err
	}

	resp := make([]shared.Model, 0)
	for _, tag := range tags {
		resp = append(resp, shared.Model{
			Name: tag.Name,
			Metadata: map[string]any{
				"size":   tag.Size,
				"digest": tag.Digest,
			},
		})
	}
	return &resp, nil
}

func (p *OllamaProvider) Chat(input shared.ChatRequest, handler shared.ProviderToolHandler) (*shared.ChatResponse, error) {
	tools := make([]api.ToolDefinition, 0)
	for _, tool := range handler.GetTools() {

		properties := make(map[string]api.ToolProperty, 0)
		for name, property := range tool.Properties {
			properties[name] = api.ToolProperty{
				Type:        string(property.Type),
				Description: property.Description,
				Enum:        property.Enum,
			}
		}

		tools = append(tools, api.ToolDefinition{
			Type: "function",
			Function: api.ToolFunction{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters: api.ToolParameters{
					Type:       string(tool.Type),
					Required:   tool.Required,
					Properties: properties,
				},
			},
		})
	}

	request := api.ChatRequest{
		Model: input.Model,
		Messages: []api.ChatMessage{
			{
				Role:    "user",
				Content: input.Message.Content,
			},
		},
		Tools:       tools,
		Stream:      false,
		KeepAlive:   -1,
		ContextSize: 4096,
	}

	var result api.ChatResponse
	hand := func(response api.ChatResponse) error {
		result = response
		return nil
	}

	if err := p.Client.Chat(p.Context, request, func(response api.ChatResponse) error {
		request.Messages = append(request.Messages, response.Message)
		// Check if we have an uncommon response without any toolcalls:
		if len(response.Message.ToolCalls) > 0 {
			// Only add the system prompt, if we have a tool response
			request.Messages = append([]api.ChatMessage{
				{
					Role: "system",
					Content: `You are a helpful assistant with tool calling capabilities.
					When you receive a tool call response, use the output to format an answer to the orginal user question.`,
				},
			}, request.Messages...)

			for _, toolcall := range response.Message.ToolCalls {
				result, err := handler.Execute(p.Context, shared.ToolFunction{
					Index:     0,
					Name:      toolcall.Function.Name,
					Arguments: toolcall.Function.Arguments,
				})
				if err != nil {
					return err
				}
				request.Messages = append(request.Messages, api.ChatMessage{
					Role:    "tool",
					Content: result,
				})
			}
		}

		return p.Client.Chat(p.Context, request, hand)
	}); err != nil {
		return nil, err
	}

	return &shared.ChatResponse{
		Model: result.Model,
		Message: shared.Message{
			Content: result.Message.Content,
		},
	}, nil
}

func (p *OllamaProvider) Embed(input shared.EmbedRequest) (*shared.EmbedResponse, error) {
	request := api.EmbedRequest{
		Model: input.Model,
		Input: input.Message.Content,
	}

	var output shared.EmbedResponse

	if err := p.Client.Embed(p.Context, request, func(response api.EmbedResponse) error {
		output = shared.EmbedResponse{
			Model:      response.Model,
			Embeddings: response.Embeddings,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &output, nil
}

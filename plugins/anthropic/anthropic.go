package anthropic

import (
	"fmt"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/liushuangls/go-anthropic/v2/jsonschema"
	"github.com/mwantia/queueverse/pkg/plugin/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p *AnthropicProvider) GetModels() (*[]shared.Model, error) {
	return &[]shared.Model{
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

func (p *AnthropicProvider) Chat(input shared.ChatRequest, handler shared.ProviderToolHandler) (*shared.ChatResponse, error) {
	tools := make([]anthropic.ToolDefinition, 0)
	for _, tool := range handler.GetTools() {

		properties := make(map[string]jsonschema.Definition, 0)
		for name, property := range tool.Properties {
			properties[name] = jsonschema.Definition{
				Type:        jsonschema.DataType(property.Type),
				Enum:        property.Enum,
				Description: property.Description,
			}
		}

		tools = append(tools, anthropic.ToolDefinition{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: jsonschema.Definition{
				Type:       jsonschema.Object,
				Required:   tool.Required,
				Properties: properties,
			},
		})
	}

	p.Logger.Info("Tools", "len", len(handler.GetTools()))
	p.Logger.Info("Tools", "len", len(tools))

	request := anthropic.MessagesRequest{
		MaxTokens: 100,
		Model:     anthropic.Model(input.Model),
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(input.Message.Content),
		},
		Tools: tools,
	}

	max := 5
	for attempts := 0; attempts < max; attempts++ {
		response, err := p.Client.CreateMessages(p.Context, request)
		if err != nil {
			return nil, fmt.Errorf("failed to create messages: %w", err)
		}

		request.Messages = append(request.Messages, anthropic.Message{
			Role:    anthropic.RoleAssistant,
			Content: response.Content,
		})

		var use *anthropic.MessageContentToolUse
		for _, content := range response.Content {
			if content.Type == anthropic.MessagesContentTypeToolUse {
				use = content.MessageContentToolUse
				break
			}
		}

		if use == nil {
			return &shared.ChatResponse{
				Model: string(response.Model),
				Message: shared.Message{
					Content: response.Content[0].GetText(),
				},
			}, nil
		}

		function := shared.ToolFunction{
			Index: 0,
			Name:  use.Name,
		}
		if err := use.UnmarshalInput(&function.Arguments); err != nil {
			return nil, err
		}

		result, err := handler.Execute(p.Context, function)
		if err != nil {
			return nil, err
		}

		request.Messages = append(request.Messages, anthropic.NewToolResultsMessage(use.ID, result, false))
	}

	return nil, fmt.Errorf("max attempts reached")
}

func (*AnthropicProvider) Embed(shared.EmbedRequest) (*shared.EmbedResponse, error) {
	return nil, status.Error(codes.Unavailable, "Embed models are not supported")
}

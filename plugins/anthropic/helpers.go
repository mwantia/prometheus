package anthropic

import (
	"fmt"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/liushuangls/go-anthropic/v2/jsonschema"

	"github.com/mwantia/queueverse/pkg/plugin/provider"
)

func CreateMessageRequest(req provider.ChatRequest) anthropic.MessagesRequest {
	messages := []anthropic.Message{}
	for _, msg := range req.Messages {
		messages = append(messages, anthropic.Message{
			Role: anthropic.ChatRole(msg.Role),
			Content: []anthropic.MessageContent{
				anthropic.NewTextMessageContent(msg.Content),
			},
		})
	}

	tools := []anthropic.ToolDefinition{}
	for _, tool := range req.Tools {
		properties := map[string]jsonschema.Definition{}
		for name, property := range tool.Parameters.Properties {
			properties[name] = jsonschema.Definition{
				Type:        jsonschema.DataType(property.Type),
				Description: property.Description,
				Enum:        property.Enum,
			}
		}

		definition := anthropic.ToolDefinition{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: jsonschema.Definition{
				Type:       jsonschema.Object,
				Properties: properties,
				Required:   tool.Parameters.Required,
			},
		}

		tools = append(tools, definition)
	}

	return anthropic.MessagesRequest{
		Model:     anthropic.Model(req.Model),
		Tools:     tools,
		Messages:  messages,
		Metadata:  req.Metadata,
		MaxTokens: 1000,
	}
}

func CreateMessageResponse(resp anthropic.MessagesResponse) (*provider.ChatResponse, error) {
	if resp.Type == anthropic.MessagesResponseTypeError {
		return nil, fmt.Errorf("error during message creation")
	}

	result := provider.ChatResponse{
		Model:    string(resp.Model),
		Messages: make([]provider.ChatMessage, 0),
	}

	for _, content := range resp.Content {
		switch content.Type {
		case anthropic.MessagesContentTypeText:
			result.Messages = append(result.Messages, provider.ChatMessage{
				Role:    provider.ChatRoleAssistant,
				Content: content.GetText(),
			})
		case anthropic.MessagesContentTypeToolUse:
			tool := content.MessageContentToolUse

			var arguments map[string]any
			if err := tool.UnmarshalInput(&arguments); err != nil {
				return nil, fmt.Errorf("failed to unmarshal tool use inputs: %w", err)
			}

			result.Messages = append(result.Messages, provider.ChatMessage{
				ID:      tool.ID,
				Role:    provider.ChatRoleTool,
				Content: content.GetText(),
				ToolCalls: []provider.ToolCall{
					{
						Function: provider.ToolFunction{
							Name:      tool.Name,
							Arguments: arguments,
						},
					},
				},
			})
		case anthropic.MessagesContentTypeImage:
			result.Messages = append(result.Messages, provider.ChatMessage{
				//	ID:      content.ID,
				Role:    provider.ChatRoleImage,
				Content: content.GetText(),
			})
		case anthropic.MessagesContentTypeDocument:
			result.Messages = append(result.Messages, provider.ChatMessage{
				//	ID:      content.ID,
				Role:    provider.ChatRoleDocument,
				Content: content.GetText(),
			})
		}
	}

	return &result, nil
}

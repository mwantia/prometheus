package anthropic

import (
	"fmt"

	"github.com/liushuangls/go-anthropic/v2"

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

	return anthropic.MessagesRequest{
		Model:    anthropic.Model(req.Model),
		Messages: messages,
		Metadata: req.Metadata,
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
				ID:      content.ID,
				Role:    provider.ChatRoleAssistant,
				Content: content.GetText(),
			})
		case anthropic.MessagesContentTypeToolUse:
			result.Messages = append(result.Messages, provider.ChatMessage{
				ID:      content.ID,
				Role:    provider.ChatRoleTool,
				Content: content.GetText(),
			})
		case anthropic.MessagesContentTypeImage:
			result.Messages = append(result.Messages, provider.ChatMessage{
				ID:      content.ID,
				Role:    provider.ChatRoleImage,
				Content: content.GetText(),
			})
		case anthropic.MessagesContentTypeDocument:
			result.Messages = append(result.Messages, provider.ChatMessage{
				ID:      content.ID,
				Role:    provider.ChatRoleDocument,
				Content: content.GetText(),
			})
		}
	}

	return &result, nil
}

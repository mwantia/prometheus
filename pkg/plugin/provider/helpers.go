package provider

import "strings"

func SystemMessage(content string) ChatMessage {
	return ChatMessage{
		Role: ChatRoleSystem,
		Content: []ChatMessageContent{
			{
				Type: ChatMessageText,
				Text: content,
			},
		},
	}
}

func UserMessage(content string) ChatMessage {
	return ChatMessage{
		Role: ChatRoleUser,
		Content: []ChatMessageContent{
			{
				Type: ChatMessageText,
				Text: content,
			},
		},
	}
}

func AssistantMessage(content string) ChatMessage {
	return ChatMessage{
		Role: ChatRoleAssistant,
		Content: []ChatMessageContent{
			{
				Type: ChatMessageText,
				Text: content,
			},
		},
	}
}

func ToolCallsMessage(content string, toolcalls []ToolCall) ChatMessage {
	return ChatMessage{
		Role: ChatRoleAssistant,
		Content: []ChatMessageContent{
			{
				Type: ChatMessageText,
				Text: content,
			},
			{
				Type:      ChatMessageToolUse,
				ToolCalls: toolcalls,
			},
		},
	}
}

func (msg *ChatMessage) GetText() string {
	var text strings.Builder
	for _, content := range msg.Content {
		if content.Type == ChatMessageText {
			text.WriteString(content.Text)
		}
	}

	return text.String()
}

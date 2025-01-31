package provider

func SystemMessage(content string) ChatMessage {
	return ChatMessage{
		Role:      ChatRoleSystem,
		Content:   content,
		ToolCalls: nil,
	}
}

func UserMessage(content string) ChatMessage {
	return ChatMessage{
		Role:      ChatRoleUser,
		Content:   content,
		ToolCalls: nil,
	}
}

func AssistantMessage(content string) ChatMessage {
	return ChatMessage{
		Role:      ChatRoleAssistant,
		Content:   content,
		ToolCalls: nil,
	}
}

func ToolMessage(content string) ChatMessage {
	return ChatMessage{
		Role:      ChatRoleTool,
		Content:   content,
		ToolCalls: nil,
	}
}

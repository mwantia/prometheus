package convert

import (
	"github.com/liushuangls/go-anthropic/v2"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
)

func ConvertMessages(input []provider.ChatMessage) []anthropic.Message {
	output := []anthropic.Message{}
	for _, message := range input {
		output = append(output, ConvertMessage(message))
	}

	return output
}

func ConvertMessage(input provider.ChatMessage) anthropic.Message {
	output := anthropic.Message{
		Role:    anthropic.ChatRole(input.Role),
		Content: nil,
	}

	return output
}

func ConvertContents(input []provider.ChatMessageContent) []anthropic.MessageContent {
	output := []anthropic.MessageContent{}
	for _, content := range input {
		output = append(output, ConvertContent(content))
	}

	return output
}

func ConvertContent(input provider.ChatMessageContent) anthropic.MessageContent {
	output := anthropic.MessageContent{}

	return output
}

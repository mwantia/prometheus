package convert

import (
	"github.com/liushuangls/go-anthropic/v2"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
)

func Convert(input provider.ChatRequest) (anthropic.MessagesRequest, error) {
	output := anthropic.MessagesRequest{
		Model: anthropic.Model(input.Model),
	}

	return output, nil
}

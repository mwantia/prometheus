package tasks

import (
	"github.com/mwantia/queueverse/pkg/ollama"
	"github.com/mwantia/queueverse/pkg/plugin/tools"
)

func convertToolService(service tools.ToolPlugin) {
}

func createTools() []ollama.Tool {
	return []ollama.Tool{
		{
			Type: "function",
			Function: ollama.ToolFunction{
				Name:        "get_current_time",
				Description: "Get the current time in the specified timezone",
				Parameters: ollama.ToolFunctionParameter{
					Type:     "string",
					Required: []string{"timezone"},
					Properties: map[string]ollama.ToolFunctionProperty{
						"timezone": {
							Type:        "string",
							Description: "The timezone to use. Must be a IANA Time Zone",
						},
					},
				},
			},
		},
	}
}

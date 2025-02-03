package tools

import "github.com/mwantia/queueverse/pkg/plugin/provider"

var TimeGetCurrent = provider.ToolDefinition{
	Name: "time_get_current",
	Description: `Get the current time in the specified timezone.
	The timezone must be a IANA compatible timezone.
	The output is in the following format 'Mon Jan 2 15:04:05'.
	Only use the toll, if the conversation specifically requires the current time.`,
	Parameters: provider.ToolParameters{
		Type:     provider.ToolTypeString,
		Required: []string{"timezone"},
		Properties: map[string]provider.ToolProperty{
			"timezone": {
				Type:        provider.ToolTypeString,
				Description: "The timezone to use. Must be a IANA compatible time zone.",
			},
		},
	},
}

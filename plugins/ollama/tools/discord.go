package tools

import "github.com/mwantia/queueverse/pkg/plugin/provider"

var DiscordSendPM = provider.ToolDefinition{
	Name:     "discord_send_pm",
	Type:     provider.TypeBoolean,
	Required: []string{"user", "message"},
	Description: `Sends a private message over Discord to a user.
	The message property supports markdown to allow, for example bold or cursiv text.
	It should only be used when the user wants to send a message, or even a tool response to another user.
	You must define the intend of the message, as well as the user it was send or requested from.
	Use templates like '{{ .displayname }}' to refer to the user making the request (e.g. Displayname = 'Max Mustermann').`,
	Properties: map[string]provider.ToolProperty{
		"user": {
			Type:        provider.TypeString,
			Description: "The user the private message will be send to.",
		},
		"message": {
			Type:        provider.TypeString,
			Description: "The private message to send over Discord; Supports markdown language.",
		},
	},
}

var DiscordListContact = provider.ToolDefinition{
	Name:     "discord_list_contact",
	Type:     provider.TypeObject,
	Required: []string{"query"},
	Description: `Retrieves a list of usernames available within Discord.
	Can be used in combination with other Discord tools that require information about a user.
	The received userdata is stored and provided to other tool calls as variables. 
	This can be accessed in a template format by defining '{{ user.<property> }}'.
	These can even be used in property values for other tool calls.
	The following variables will become available after searching for a user:
	 * user.username
	 * user.displayname
	 * user.mail
	 * user.status`,
	Properties: map[string]provider.ToolProperty{
		"query": {
			Type: provider.TypeString,
			Description: `Defines the search query used to find the correct user.
			 This can be the displayname, surname, lastname or other available contact information`,
		},
	},
}

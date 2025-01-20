package ollama

type Style string

const (
	ConciseStyle Style = "concise"
	DefaultStyle Style = "default"
	FormalStyle  Style = "formal"
)

type RespondStyle struct {
	Style    Style  `json:"style"`
	Template string `json:"template"`
}

var Styles = map[string]RespondStyle{
	string(ConciseStyle): {
		Style:    ConciseStyle,
		Template: "Reply with shorter responses and keep the messages concise and to the point.",
	},
	string(DefaultStyle): {
		Style:    DefaultStyle,
		Template: "",
	},
	string(FormalStyle): {
		Style:    FormalStyle,
		Template: "Reply in a clear, well-structured and formal tone.",
	},
}

func (c *Client) addSystemStylePrompt(req *ChatRequest, data any) error {
	n := string(c.Style)
	style, exists := Styles[n]
	if exists && style.Template != "" {
		return c.AddSystemPrompt(req, style.Template, data)
	}

	return nil
}

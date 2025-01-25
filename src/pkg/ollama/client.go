package ollama

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"
)

type Client struct {
	http *http.Client `json:"-"`

	Endpoint string `json:"endpoint"`
	Model    string `json:"model,omitempty"`
	Style    Style  `json:"style,omitempty"`
}

func CreateClient(endpoint, model string, http *http.Client) *Client {
	return &Client{
		http:     http,
		Endpoint: endpoint,
		Model:    model,
		Style:    DefaultStyle,
	}
}

func (c *Client) AddSystemPrompt(req *ChatRequest, system string, data any) error {
	tmpl, err := template.New("system_prompt").Parse(system)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}

	req.Messages = append([]ChatMessage{{
		Role:    "system",
		Content: buf.String(),
	}}, req.Messages...)

	return nil
}

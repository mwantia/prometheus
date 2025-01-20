package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Parameters  ToolFunctionParameter `json:"parameters"`
}

type ToolFunctionParameter struct {
	Type       string                          `json:"type"`
	Required   []string                        `json:"required"`
	Properties map[string]ToolFunctionProperty `json:"properties,omitempty"`
}

type ToolFunctionProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

type ToolCall struct {
	Function ToolCallFunction `json:"function"`
}

type ToolCallFunction struct {
	Index     int            `json:"index,omitempty"`
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments,omitempty"`
}

const ToolSystemPrompt = `You are a helpful assistant with tool calling capabilities.
When you receive a tool call response, use the output to format an answer to the orginal user question.`

func (tc *ToolCall) Debug() {
	log.Println("----- Tool Call Function -----")
	log.Printf("Name: %s", tc.Function.Name)
	for k, v := range tc.Function.Arguments {
		log.Printf("Argument '%s': %v", k, v)
	}
	log.Println("------------------------------")
}

func (c *Client) updateSystemTool(req *ChatRequest, data any) error {
	tmpl, err := template.New("system_tool").Parse(ToolSystemPrompt)
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

func (c *Client) ChatTools(ctx context.Context, req ChatRequest, res ChatResponseHandler, tools []Tool) error {
	if err := c.stream(ctx, http.MethodPost, "/api/chat", struct {
		Tools []Tool `json:"tools,omitempty"`
		ChatRequest
	}{Tools: tools, ChatRequest: req}, func(bts []byte) error {
		// Received the tool response from Ollama
		var resp ChatResponse
		if err := json.Unmarshal(bts, &resp); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
		// Add the model's response to the conversation history
		req.Messages = append(req.Messages, resp.Message)
		// Check if we have an uncommon response without any toolcalls:
		if len(resp.Message.ToolCalls) > 0 {
			// Only add the system prompt, if we have a tool response
			if err := c.AddSystemPrompt(&req, ToolSystemPrompt, struct{}{}); err != nil {
				return fmt.Errorf("system prompt error: %w", err)
			}
			// Add our new message, including the result as content
			for _, toolcall := range resp.Message.ToolCalls {
				if toolcall.Function.Name == "get_current_time" {
					tz := toolcall.Function.Arguments["timezone"]
					location, _ := time.LoadLocation(tz.(string))

					req.Messages = append(req.Messages, ChatMessage{
						Role:    "tool",
						Content: time.Now().In(location).Format("Mon Jan 2 15:04:05"),
					})
				}
			}

			return c.Chat(ctx, req, res)
		}
		// Only bypass tool response, if we receive an uncommon reply with content
		if resp.Message.Content != "" {
			return res(resp)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

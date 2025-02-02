package api

import "time"

const (
	UserAgent        = "QueueVerse/OllamaClient"
	MaxScannerBuffer = 512 * 1000
)

type DataHandler func([]byte) error

type Tag struct {
	Name   string `json:"name"`
	Size   int    `json:"size"`
	Digest string `json:"digest"`
}

type ChatRequest struct {
	Model       string           `json:"model"`
	Messages    []ChatMessage    `json:"messages"`
	Tools       []ToolDefinition `json:"tools,omitempty"`
	Stream      bool             `json:"stream,omitempty"`
	KeepAlive   int              `json:"keep_alive,omitempty"`
	ContextSize int              `json:"context_size,omitempty"`
}

type ChatMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ChatResponse struct {
	Model      string      `json:"model"`
	CreatedAt  time.Time   `json:"created_at"`
	Message    ChatMessage `json:"message"`
	DoneReason string      `json:"done_reason,omitempty"`
	Done       bool        `json:"done"`
}

type ChatResponseHandler func(ChatResponse) error

type EmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbedResponse struct {
	Model      string      `json:"model"`
	Embeddings [][]float32 `json:"embeddings"`
}

type EmbedResponseHandler func(EmbedResponse) error

type ToolDefinition struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  ToolParameters `json:"parameters"`
}

type ToolParameters struct {
	Type       string                  `json:"type"`
	Required   []string                `json:"required"`
	Properties map[string]ToolProperty `json:"properties,omitempty"`
}

type ToolProperty struct {
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

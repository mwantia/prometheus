package api

import "time"

const (
	UserAgent        = "QueueVerse/OllamaClient"
	MaxScannerBuffer = 512 * 1000
)

type DataHandler func([]byte) error

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	KeepAlive   int           `json:"keep_alive"`
	ContextSize int           `json:"context_size"`
}

type ChatResponse struct {
	Model      string      `json:"model"`
	CreatedAt  time.Time   `json:"created_at"`
	Message    ChatMessage `json:"message"`
	DoneReason string      `json:"done_reason,omitempty"`

	Done bool `json:"done"`
}

type ChatResponseHandler func(ChatResponse) error

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TagModel struct {
	Name   string `json:"name"`
	Size   int    `json:"size"`
	Digest string `json:"digest"`
}

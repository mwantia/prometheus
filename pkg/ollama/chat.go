package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

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
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	Role      string     `json:"role"`
	Content   string     `json:"content"`
}

func (c *Client) Chat(ctx context.Context, req ChatRequest, res ChatResponseHandler) error {
	req.ContextSize = 8192
	req.KeepAlive = -1

	if err := c.addSystemStylePrompt(&req, struct{}{}); err != nil {
		return fmt.Errorf("system style prompt error: %w", err)
	}
	return c.stream(ctx, http.MethodPost, "/api/chat", req, func(bts []byte) error {
		var resp ChatResponse
		if err := json.Unmarshal(bts, &resp); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
		return res(resp)
	})
}

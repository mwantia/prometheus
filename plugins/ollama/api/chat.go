package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *clientImpl) Chat(ctx context.Context, req ChatRequest, handler ChatResponseHandler) error {
	return c.stream(ctx, http.MethodPost, "/api/chat", req, func(data []byte) error {
		var resp ChatResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
		return handler(resp)
	})
}

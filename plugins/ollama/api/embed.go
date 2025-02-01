package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *clientImpl) Embed(ctx context.Context, req EmbedRequest, handler EmbedResponseHandler) error {
	return c.do(ctx, http.MethodPost, "/api/embed", req, func(data []byte) error {
		var result EmbedResponse
		if err := json.Unmarshal(data, &result); err != nil {
			return fmt.Errorf("unmarschal: %w", err)
		}
		return handler(result)
	})
}

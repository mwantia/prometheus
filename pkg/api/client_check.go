package api

import (
	"context"
	"fmt"
	"net/http"
)

type CheckHandler func(int) error

func (c *clientImpl) check(ctx context.Context, path string, handler CheckHandler) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, c.config.Endpoint+path, nil)
	if err != nil {
		return fmt.Errorf("unable to create http request: %w", err)
	}

	req.Header.Set("User-Agent", useragent)
	if c.config.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.Token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("error during http request: %w", err)
	}
	if resp == nil {
		return fmt.Errorf("nil response received")
	}

	defer resp.Body.Close()

	return handler(resp.StatusCode)
}

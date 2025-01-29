package api

import (
	"context"
	"fmt"
	"net/http"
)

func (c *clientImpl) Health(ctx context.Context) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, fmt.Errorf("context error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, c.config.Endpoint, nil)
	if err != nil {
		return false, fmt.Errorf("unable to create http request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return false, fmt.Errorf("error during http request: %w", err)
	}
	if resp == nil {
		return false, fmt.Errorf("nil response received")
	}

	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

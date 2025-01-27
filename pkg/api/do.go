package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DoHandler func([]byte) error

func (c *clientImpl) do(ctx context.Context, method, path string, data any, handler DoHandler) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	buf, err := createBuffer(data)
	if err != nil {
		return fmt.Errorf("failed to create buffer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.config.Endpoint+path, buf)
	if err != nil {
		return fmt.Errorf("unable to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
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

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("request failed with status '%d'", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	return handler(body)
}

func createBuffer(data any) (*bytes.Buffer, error) {
	if data == nil {
		return &bytes.Buffer{}, nil
	}

	buf, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %w", err)
	}

	return bytes.NewBuffer(buf), nil
}

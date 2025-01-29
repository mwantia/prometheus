package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *clientImpl) do(ctx context.Context, method, path string, data any, handler DataHandler) error {
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
	req.Header.Set("User-Agent", UserAgent)

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

func (c *clientImpl) stream(ctx context.Context, method, path string, data any, handler DataHandler) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	buf, err := createBuffer(data)
	if err != nil {
		return fmt.Errorf("create buffer error: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, method, c.config.Endpoint+path, buf)
	if err != nil {
		return fmt.Errorf("unable to create http request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/x-ndjson")
	request.Header.Set("User-Agent", UserAgent)

	response, err := c.http.Do(request)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	if response == nil {
		return fmt.Errorf("nil response received")
	}

	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("request failed with status %d", response.StatusCode)
	}

	scanner := bufio.NewScanner(response.Body)
	scanBuf := make([]byte, 0, MaxScannerBuffer)
	scanner.Buffer(scanBuf, MaxScannerBuffer)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled during scan: %w", ctx.Err())
		default:
			var result struct {
				Error string `json:"error,omitempty"`
			}

			bts := scanner.Bytes()
			if err := json.Unmarshal(bts, &result); err != nil {
				return fmt.Errorf("unmarshal: %w", err)
			}

			if result.Error != "" {
				return fmt.Errorf("response error: %s", result.Error)
			}

			if err := handler(bts); err != nil {
				return fmt.Errorf("handler error: %w", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
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

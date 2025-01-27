package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type EmbedRequest struct {
	Model string `json:"model,omitempty"`
	Input any    `json:"input"`
}

type EmbedResponse struct {
	Model      string      `json:"model"`
	Embeddings [][]float32 `json:"embeddings"`
}

type EmbedResponseHandler func(EmbedResponse) error

func (c *Client) Embed(ctx context.Context, req EmbedRequest) (*EmbedResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	if req.Model == "" {
		req.Model = c.Model
	}

	buf, err := createBuffer(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create buffer: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Endpoint+"/api/embed", buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/x-ndjson")
	request.Header.Set("User-Agent", UserAgent)

	response, err := c.http.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error during http request: %w", err)
	}
	if response == nil {
		return nil, fmt.Errorf("nil response received")
	}

	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("request failed with status '%d'", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var resp *EmbedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %w", err)
	}

	return resp, nil
}

func (c *Client) EmbedHandle(ctx context.Context, req EmbedRequest, handler EmbedResponseHandler) error {
	resp, err := c.Embed(ctx, req)
	if err != nil {
		return err
	}

	return handler(*resp)
}

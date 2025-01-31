package api

import (
	"context"
	"encoding/json"
	"net/http"
)

func (c *clientImpl) Tags(ctx context.Context) ([]Tag, error) {
	type TagsResponse struct {
		Models []Tag `json:"models"`
	}
	var resp TagsResponse

	if err := c.do(ctx, http.MethodGet, "/api/tags", nil, func(data []byte) error {
		if err := json.Unmarshal(data, &resp); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return resp.Models, nil
}

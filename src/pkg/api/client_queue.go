package api

import (
	"context"
	"encoding/json"
	"net/http"
)

func (c *clientImpl) QueueTask(ctx context.Context, req QueueRequest) (QueueResponse, error) {
	var resp QueueResponse
	if err := c.do(ctx, http.MethodPost, "/v1/queue", req, func(body []byte) error {
		return json.Unmarshal(body, &resp)
	}); err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *clientImpl) QueueTaskResult(ctx context.Context, task string) (QueueResult, error) {
	var result QueueResult
	if err := c.do(ctx, http.MethodGet, "/v1/queue/"+task, nil, func(body []byte) error {
		var resp QueueResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return err
		}

		result = resp.Result

		return nil
	}); err != nil {
		return result, err
	}

	return result, nil
}

func (c *clientImpl) ListQueuedTasks(ctx context.Context) ([]QueueResult, error) {
	var results []QueueResult
	if err := c.do(ctx, http.MethodGet, "/v1/queue", nil, func(body []byte) error {
		var resp []QueueResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return err
		}

		for _, r := range resp {
			results = append(results, r.Result)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return results, nil
}

func (c *clientImpl) QueueState(ctx context.Context, task string) (bool, error) {
	var result bool
	if err := c.check(ctx, "/v1/queue/"+task, func(statuscode int) error {
		if statuscode == http.StatusOK {
			result = true
		}
		return nil
	}); err != nil {
		return false, err
	}

	return result, nil
}

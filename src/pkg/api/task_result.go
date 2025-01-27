package api

import (
	"context"
	"fmt"
)

func (t *taskImpl) Result(ctx context.Context) (string, error) {
	if t.task == "" {
		return "", fmt.Errorf("no task running in queue")
	}

	result, err := t.client.QueueTaskResult(ctx, t.task)
	if err != nil {
		return "", err
	}

	return result.Text, nil
}

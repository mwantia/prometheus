package api

import (
	"context"
)

func (t *taskImpl) Queue(ctx context.Context, prompt string) error {
	resp, err := t.client.QueueTask(ctx, QueueRequest{
		Prompt: prompt,
		Model:  t.config.Model,
		Style:  t.config.Style,
	})
	if err != nil {
		return err
	}

	t.task = resp.Task
	return nil
}

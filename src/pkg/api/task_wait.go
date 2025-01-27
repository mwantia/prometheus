package api

import (
	"context"
	"fmt"
	"time"
)

func (t *taskImpl) Wait(ctx context.Context) error {
	const base = 100 * time.Millisecond
	const max = 2 * time.Second

	delay := base

	if t.task == "" {
		return fmt.Errorf("no task running in queue")
	}

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		done, err := t.client.QueueState(ctx, t.task)
		if err != nil {
			return fmt.Errorf("error checking queue state: %w", err)
		}

		if done {
			break
		}

		select {
		case <-time.After(delay):
			if delay < max {
				delay *= 2
				if delay > max {
					delay = max
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

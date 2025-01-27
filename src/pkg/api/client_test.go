package api

import (
	"context"
	"net/http"
	"os"
	"testing"
)

func TestClient(tst *testing.T) {
	ctx := context.TODO()
	client := NewClient(http.DefaultClient, ClientConfig{
		Endpoint: os.Getenv("Endpoint"),
		Token:    os.Getenv("TOKEN"),
	})

	tst.Run("Test.ListQueuedTasks", func(t *testing.T) {
		results, err := client.ListQueuedTasks(ctx)
		if err != nil {
			t.Errorf("Failed to list queued tasks: %v", err)
		}

		for _, result := range results {
			t.Logf("Result: %s", result.Text)
		}
	})
}

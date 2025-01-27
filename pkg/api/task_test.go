package api

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	ctx := context.TODO()
	task := NewTask(http.DefaultClient, TaskConfig{
		Endpoint: "http://127.0.0.1:8080",
		Token:    "",
	})

	t.Run("Test.Task", func(t *testing.T) {
		if err := task.Queue(ctx, "Tell me the current time in germany"); err != nil {
			t.Error(err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := task.Wait(ctx); err != nil {
			t.Error(err.Error())
		}

		text, err := task.Result(ctx)
		if err != nil {
			t.Error(err.Error())
		}

		log.Print(text)
		t.Log(text)
	})
}

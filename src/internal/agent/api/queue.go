package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hibiken/asynq"
	"github.com/mwantia/prometheus/pkg/tasks"
)

type GeneratePromptRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model,omitempty"`
}

type GeneratePromptResponse struct {
	TaskID string `json:"taskid"`
	State  string `json:"state"`
	Queue  string `json:"queue"`
	Result string `json:"result,omitempty"`
}

func HandleQueueGet(w http.ResponseWriter, r *http.Request, address string, db int) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")

	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr: address,
		DB:   db,
	})
	defer inspector.Close()

	queue := r.URL.Query().Get("queue")
	if queue == "" {
		queue = "default"
	}
	taskid := r.URL.Query().Get("taskid")
	if taskid == "" {
		infos, err := inspector.ListCompletedTasks(queue)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Unable to get info for task: %v", err)
			return
		}

		res := make([]GeneratePromptResponse, 0)
		for _, info := range infos {
			res = append(res, GeneratePromptResponse{
				TaskID: info.ID,
				State:  info.State.String(),
				Queue:  info.Queue,
				Result: string(info.Result),
			})
		}

		if err := encoder.Encode(res); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Unable to encode final result: %v", err)
			return
		}
	} else {
		info, err := inspector.GetTaskInfo(queue, taskid)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Unable to get info for task: %v", err)
			return
		}

		res := GeneratePromptResponse{
			TaskID: info.ID,
			State:  info.State.String(),
			Queue:  info.Queue,
			Result: string(info.Result),
		}

		if err := encoder.Encode(res); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Unable to encode final result: %v", err)
			return
		}
	}
}

func HandleQueuePost(w http.ResponseWriter, r *http.Request, address string, db int) {
	var request GeneratePromptRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unable to decode body: %v", err)
		return
	}

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: address,
		DB:   db,
	})
	defer client.Close()

	prompt, err := json.Marshal(tasks.GeneratePrompt{
		Content: request.Prompt,
		Model:   request.Model,
	})
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Unable to create task: %v", err)
		return
	}

	task := asynq.NewTask(tasks.TaskTypeGeneratePrompt, prompt)
	taskid := fmt.Sprintf("T%d", time.Now().UnixNano())

	log.Printf("Task: %s", taskid)

	info, err := client.EnqueueContext(r.Context(), task, asynq.TaskID(taskid), asynq.Retention(7*24*time.Hour))
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Unable to enqueue task: %v", err)
		return
	}

	log.Printf("Enqueued task: id=%s queue=%s", info.ID, info.Queue)

	response := GeneratePromptResponse{
		TaskID: info.ID,
		State:  info.State.String(),
		Queue:  info.Queue,
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(response); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Unable to encode response: %v", err)
		return
	}
}

func HandleQueue(address string, db int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			HandleQueueGet(w, r, address, db)

		case http.MethodPost:
			HandleQueuePost(w, r, address, db)
		}
	}
}

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/mwantia/prometheus/internal/config"
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

func HandleGetQueue(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		inspector := asynq.NewInspector(asynq.RedisClientOpt{
			Addr:     cfg.Redis.Endpoint,
			DB:       cfg.Redis.Database,
			Password: cfg.Redis.Password,
		})
		defer inspector.Close()

		queue := c.DefaultQuery("queue", "default")
		taskid := c.Query("taskid")

		if taskid == "" {
			infos, err := inspector.ListCompletedTasks(queue)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
					"error": fmt.Sprintf("Unable to get task info for queue '%s': %v", queue, err),
				})
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

			c.JSON(http.StatusOK, res)
			return
		}

		info, err := inspector.GetTaskInfo(queue, taskid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": fmt.Sprintf("Unable to get task info '%s': %v", queue, err),
			})
			return
		}

		res := GeneratePromptResponse{
			TaskID: info.ID,
			State:  info.State.String(),
			Queue:  info.Queue,
			Result: string(info.Result),
		}
		c.JSON(http.StatusOK, res)
	}
}

func HandlePostQueue(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request GeneratePromptRequest
		if err := c.BindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Unable to decode body: %v", err),
			})
			return
		}

		client := asynq.NewClient(asynq.RedisClientOpt{
			Addr:     cfg.Redis.Endpoint,
			DB:       cfg.Redis.Database,
			Password: cfg.Redis.Password,
		})
		defer client.Close()

		prompt, err := json.Marshal(tasks.GeneratePrompt{
			Content: request.Prompt,
			Model:   request.Model,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": fmt.Sprintf("Unable to create task: %v", err),
			})
			return
		}

		task := asynq.NewTask(tasks.TaskTypeGeneratePrompt, prompt)
		taskid := fmt.Sprintf("T%d", time.Now().UnixNano())

		log.Printf("Task: %s", taskid)

		info, err := client.EnqueueContext(c.Request.Context(), task, asynq.TaskID(taskid), asynq.Retention(7*24*time.Hour))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": fmt.Sprintf("Unable to enqueue task: %v", err),
			})
			return
		}

		log.Printf("Enqueued task: id=%s queue=%s", info.ID, info.Queue)

		c.JSON(http.StatusOK, GeneratePromptResponse{
			TaskID: info.ID,
			State:  info.State.String(),
			Queue:  info.Queue,
		})
	}
}

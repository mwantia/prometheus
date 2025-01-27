package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/mwantia/queueverse/internal/config"
	"github.com/mwantia/queueverse/pkg/tasks"
)

func HandleGetQueueTask(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		inspector := asynq.NewInspector(asynq.RedisClientOpt{
			Addr:     cfg.Redis.Endpoint,
			DB:       cfg.Redis.Database,
			Password: cfg.Redis.Password,
		})
		defer inspector.Close()

		task := c.Param("task")

		info, err := inspector.GetTaskInfo(cfg.PoolName, task)
		if err != nil {
			switch err.Error() {
			case "asynq: queue not found":
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"pool":  cfg.PoolName,
					"error": "Queue not found",
				})
			case "asynq: task not found":
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"pool":  cfg.PoolName,
					"task":  task,
					"error": "Task not found",
				})
			default:
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("Unable to get task info: %v", err),
				})
			}

			return
		}

		res, err := tasks.CreateGenerateResponse(info)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Unable to create response: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func HandleIsQueueTaskDone(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		inspector := asynq.NewInspector(asynq.RedisClientOpt{
			Addr:     cfg.Redis.Endpoint,
			DB:       cfg.Redis.Database,
			Password: cfg.Redis.Password,
		})
		defer inspector.Close()

		task := c.Param("task")

		info, err := inspector.GetTaskInfo(cfg.PoolName, task)
		if err != nil {
			switch err.Error() {
			case "asynq: queue not found":
				c.Status(http.StatusNotFound)
			case "asynq: task not found":
				c.Status(http.StatusNotFound)
			default:
				c.Status(http.StatusBadRequest)
			}

			return
		}

		c.Status(func() int {
			switch info.State {
			case asynq.TaskStateActive, asynq.TaskStateRetry, asynq.TaskStatePending, asynq.TaskStateScheduled:
				return http.StatusAccepted
			case asynq.TaskStateArchived:
				return http.StatusGone
			case asynq.TaskStateCompleted:
				return http.StatusOK
			default:
				return http.StatusBadRequest
			}
		}())
	}
}

func HandleGetQueue(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		inspector := asynq.NewInspector(asynq.RedisClientOpt{
			Addr:     cfg.Redis.Endpoint,
			DB:       cfg.Redis.Database,
			Password: cfg.Redis.Password,
		})
		defer inspector.Close()

		infos, err := inspector.ListCompletedTasks(cfg.PoolName)
		if err != nil {
			switch err.Error() {
			case "asynq: queue not found":
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"pool":  cfg.PoolName,
					"error": "Pool not found",
				})
			default:
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("Unable to get task info: %v", err),
				})
			}

			return
		}

		results := make([]tasks.GenerateResponse, 0)
		for _, info := range infos {
			res, err := tasks.CreateGenerateResponse(info)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("Unable to create response: %v", err),
				})
				return
			}

			results = append(results, *res)
		}

		c.JSON(http.StatusOK, results)
	}
}

func HandlePostQueue(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req tasks.GenerateRequest

		if err := c.BindJSON(&req); err != nil {
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

		t, err := tasks.CreateGenerateTask(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Unable to create task: %v", err),
			})
		}

		ctx := c.Request.Context()
		task := tasks.GenerateTaskId()

		info, err := client.EnqueueContext(ctx, t, asynq.Queue(cfg.PoolName), asynq.TaskID(task), asynq.Retention(7*24*time.Hour))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Unable to enqueue task: %v", err),
			})
			return
		}

		res, err := tasks.CreateGenerateResponse(info)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Unable to create response: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

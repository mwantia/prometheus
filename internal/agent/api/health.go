package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mwantia/queueverse/internal/registry"
)

type Health struct {
	Status  string         `json:"status"`
	Healthy bool           `json:"healthy,omitempty"`
	Plugins []PluginHealth `json:"plugins"`
}

type PluginHealth struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Healthy bool   `json:"healthy,omitempty"`
	Error   string `json:"error,omitempty"`
}

func HandleGetHealth(reg *registry.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		health := Health{
			Status:  "OK",
			Healthy: true,
		}

		for _, plugin := range reg.GetPlugins() {
			err := ""
			stat := "OK"

			if !plugin.Status.IsHealthy {
				stat = "ERROR"
				health.Healthy = false

				if plugin.Status.LastKnownError != nil {
					err = plugin.Status.LastKnownError.Error()
				}
			}

			health.Plugins = append(health.Plugins, PluginHealth{
				Name:    plugin.Info.Name,
				Status:  stat,
				Healthy: plugin.Status.IsHealthy,
				Error:   err,
			})
		}

		if !health.Healthy {
			health.Status = "ERROR"
			c.JSON(http.StatusServiceUnavailable, health)
			return
		}

		c.JSON(http.StatusOK, health)
	}
}

func HandleIsHealthy(reg *registry.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, plugin := range reg.GetPlugins() {
			if !plugin.Status.IsHealthy {
				c.Status(http.StatusServiceUnavailable)
				return
			}
		}

		c.Status(http.StatusOK)
	}
}

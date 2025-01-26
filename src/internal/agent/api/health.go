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

func HandleGetHealth(reg *registry.PluginRegistry) gin.HandlerFunc {
	return func(c *gin.Context) {
		health := Health{
			Status:  "OK",
			Healthy: true,
		}

		for _, plugin := range reg.GetPlugins() {
			err := ""
			stat := "OK"

			if !plugin.IsHealthy {
				stat = "ERROR"
				health.Healthy = false

				if plugin.LastKnownError != nil {
					err = plugin.LastKnownError.Error()
				}
			}

			health.Plugins = append(health.Plugins, PluginHealth{
				Name:    plugin.Name,
				Status:  stat,
				Healthy: plugin.IsHealthy,
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

func HandleIsHealthy(reg *registry.PluginRegistry) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, plugin := range reg.GetPlugins() {
			if !plugin.IsHealthy {
				c.Status(http.StatusServiceUnavailable)
				return
			}
		}

		c.Status(http.StatusOK)
	}
}

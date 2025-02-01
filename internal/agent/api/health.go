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
		plugins := reg.GetPlugins()

		result := []registry.RegistryStatus{}
		for _, plugin := range plugins {
			result = append(result, plugin.Status)
		}

		c.JSON(http.StatusOK, result)
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

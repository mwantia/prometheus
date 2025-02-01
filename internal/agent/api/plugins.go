package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mwantia/queueverse/internal/registry"
	"github.com/mwantia/queueverse/pkg/plugin/base"
)

func HandleGetPlugins(reg *registry.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		type Result struct {
			Info         base.PluginInfo         `json:"info"`
			Capabilities base.PluginCapabilities `json:"capabilities"`
		}

		plugins := reg.GetPlugins()

		result := []Result{}
		for _, plugin := range plugins {
			result = append(result, Result{
				Info:         plugin.Info,
				Capabilities: plugin.Capabilities,
			})
		}

		c.JSON(http.StatusOK, result)
	}
}

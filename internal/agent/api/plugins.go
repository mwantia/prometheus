package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mwantia/queueverse/internal/registry"
	"github.com/mwantia/queueverse/pkg/plugin/base"
)

func HandleGetPlugins(reg *registry.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		plugins := reg.GetPlugins()

		result := []base.PluginInfo{}
		for _, plugin := range plugins {
			result = append(result, plugin.Info)
		}

		c.JSON(http.StatusOK, result)
	}
}

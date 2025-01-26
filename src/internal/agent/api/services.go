package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mwantia/queueverse/internal/registry"
	"github.com/mwantia/queueverse/pkg/plugin/identity"
)

func HandleServices(reg *registry.PluginRegistry) gin.HandlerFunc {
	return func(c *gin.Context) {
		services := make([]identity.PluginServiceInfo, 0)

		plugins := reg.GetPlugins()
		for _, plugin := range plugins {
			info, err := plugin.Services.Identity.GetPluginInfo()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
					"error": fmt.Sprintf("Unable to enqueue task: %v", err),
				})
				return
			}

			services = append(services, info.Services...)
		}

		c.JSON(http.StatusOK, services)
	}
}

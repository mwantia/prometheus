package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mwantia/queueverse/internal/registry"
)

func HandlePlugins(reg *registry.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		// plugins := reg.GetPlugins()
		c.JSON(http.StatusOK, nil)
	}
}

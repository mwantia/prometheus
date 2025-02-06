package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mwantia/queueverse/internal/registry"
	"github.com/mwantia/queueverse/pkg/plugin/shared"
)

func HandleGetModels(reg *registry.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		providers, err := reg.GetProviders()
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		result := []shared.Model{}
		for _, prov := range providers {

			info, err := prov.GetPluginInfo()
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
			}

			models, err := prov.GetModels()
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
			}

			for _, model := range *models {
				model.Metadata["provider"] = info.Name
				result = append(result, model)
			}
		}

		c.JSON(http.StatusOK, result)
	}
}

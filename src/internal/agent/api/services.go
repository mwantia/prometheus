package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mwantia/prometheus/internal/registry"
	"github.com/mwantia/prometheus/pkg/plugin/identity"
)

func HandleServices(reg *registry.PluginRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "    ")

		services := make([]identity.PluginServiceInfo, 0)

		plugins := reg.GetPlugins()
		for _, plugin := range plugins {
			info, err := plugin.Services.Identity.GetPluginInfo()
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintf(w, "Unable to enqueue task: %v", err)
				return
			}

			services = append(services, info.Services...)
		}

		encoder.Encode(services)
	}
}

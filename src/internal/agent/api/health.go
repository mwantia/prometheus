package api

import (
	"encoding/json"
	"net/http"

	"github.com/mwantia/prometheus/internal/registry"
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

func HandleHealth(reg *registry.PluginRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

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
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "    ")

		if err := encoder.Encode(health); err != nil {
			encoder.Encode(map[string]string{
				"error": err.Error(),
			})
		}
	}
}

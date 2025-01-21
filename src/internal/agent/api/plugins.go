package api

import (
	"encoding/json"
	"net/http"

	"github.com/mwantia/prometheus/internal/registry"
)

func HandlePlugins(reg *registry.PluginRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "    ")

		plugins := reg.GetPlugins()
		encoder.Encode(plugins)
	}
}

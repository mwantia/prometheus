package agent

import (
	"encoding/json"
	"net/http"
)

type RegisterPluginRequest struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address"`
}

func (a *PrometheusAgent) handlListPlugins() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "    ")

		plugins := a.Registry.GetPlugins()

		w.WriteHeader(http.StatusOK)
		encoder.Encode(plugins)
	}
}

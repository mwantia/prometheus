package agent

import "net/http"

func (a *PrometheusAgent) startServer() (*http.Server, error) {
	mux := http.NewServeMux()
	if err := a.addRoutes(mux); err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:    a.Config.Agent.Server.Address,
		Handler: mux,
	}, nil
}

func (a *PrometheusAgent) addRoutes(mux *http.ServeMux) error {
	mux.HandleFunc("/health", a.handleHealth())
	mux.HandleFunc("/v1/plugin/list", a.handlListPlugins())
	mux.Handle("/", http.NotFoundHandler())

	return nil
}

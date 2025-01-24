package ops

import (
	"context"
	"log"
	"net/http"

	"github.com/mwantia/prometheus/internal/config"
	"github.com/mwantia/prometheus/internal/registry"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	Operation

	srv *http.Server
	mux *http.ServeMux

	Namespace string
}

func (m *Metrics) Create(cfg *config.Config, reg *registry.PluginRegistry) (Cleanup, error) {
	m.mux = http.NewServeMux()
	m.srv = &http.Server{
		Addr:    cfg.Metrics.Address,
		Handler: m.mux,
	}

	m.mux.Handle("/metrics", promhttp.Handler())

	return func(ctx context.Context) error {
		return m.srv.Shutdown(ctx)
	}, nil
}

func (m *Metrics) Serve(ctx context.Context) error {
	log.Printf("Serving metrics server: %s", m.srv.Addr)
	if err := m.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

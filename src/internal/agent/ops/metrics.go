package ops

import (
	"context"
	"fmt"
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

	if err := m.addMetricsHandler(); err != nil {
		return nil, fmt.Errorf("error adding metrics: %w", err)
	}

	return func(ctx context.Context) error {
		return m.srv.Shutdown(ctx)
	}, nil
}

func (m *Metrics) addMetricsHandler() error {
	m.mux.Handle("/metrics", promhttp.Handler())
	return nil
}

func (m *Metrics) Serve(ctx context.Context) error {
	log.Printf("Serving metrics server: %s", m.srv.Addr)
	if err := m.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

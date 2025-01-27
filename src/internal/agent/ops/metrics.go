package ops

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mwantia/queueverse/internal/config"
	"github.com/mwantia/queueverse/internal/metrics"
	"github.com/mwantia/queueverse/internal/registry"
	"github.com/mwantia/queueverse/pkg/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	Operation

	srv *http.Server
	mux *http.ServeMux

	Log       log.Logger
	Namespace string
}

func (m *Metrics) Create(cfg *config.Config, reg *registry.PluginRegistry) (Cleanup, error) {
	m.Log = log.New("metrics")
	m.mux = http.NewServeMux()
	m.srv = &http.Server{
		Addr:    cfg.Metrics.Address,
		Handler: m.mux,
	}

	prom, err := metrics.Setup(cfg.PoolName)
	if err != nil {
		return nil, fmt.Errorf("failed to setup metrics: %w", err)
	}

	m.mux.Handle("/metrics", promhttp.HandlerFor(prom, promhttp.HandlerOpts{}))

	return func(ctx context.Context) error {
		return m.srv.Shutdown(ctx)
	}, nil
}

func (m *Metrics) Serve(ctx context.Context) error {
	m.Log.Info("Serving metrics server", "addr", m.srv.Addr)
	if err := m.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

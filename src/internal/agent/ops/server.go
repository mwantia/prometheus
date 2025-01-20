package ops

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/mwantia/prometheus/internal/agent/api"
	"github.com/mwantia/prometheus/internal/config"
	"github.com/mwantia/prometheus/internal/registry"
)

type Server struct {
	Operation

	mux *http.ServeMux
	srv *http.Server
}

func (s *Server) Create(cfg *config.Config, reg *registry.PluginRegistry) (Cleanup, error) {
	s.mux = http.NewServeMux()
	if err := s.addRoutes(cfg, reg); err != nil {
		return nil, fmt.Errorf("error adding routes: %w", err)
	}

	s.srv = &http.Server{
		Addr:    cfg.Server.Address,
		Handler: s.mux,
	}

	return func(ctx context.Context) error {
		return s.srv.Shutdown(ctx)
	}, nil
}

func (s *Server) Serve(ctx context.Context) error {
	log.Printf("Serving http server: %s", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) addRoutes(cfg *config.Config, reg *registry.PluginRegistry) error {
	s.mux.HandleFunc("/v1/health", api.HandleHealth(reg))
	s.mux.HandleFunc("/v1/plugins/list", api.HandleListPlugins(reg))
	s.mux.HandleFunc("/v1/queue", api.HandleQueue(cfg.Redis.Address, cfg.Redis.Database))

	s.mux.Handle("/", http.NotFoundHandler())

	return nil
}

package ops

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/mwantia/prometheus/internal/agent/api"
	"github.com/mwantia/prometheus/internal/config"
	"github.com/mwantia/prometheus/internal/metrics"
	"github.com/mwantia/prometheus/internal/registry"
)

type Server struct {
	Operation

	engine *gin.Engine
	srv    *http.Server
}

func (s *Server) Create(cfg *config.Config, reg *registry.PluginRegistry) (Cleanup, error) {
	s.engine = gin.Default()
	s.srv = &http.Server{
		Addr:    cfg.Server.Address,
		Handler: s.engine,
	}

	if err := s.addMiddlewares(); err != nil {
		return nil, fmt.Errorf("error adding routes: %w", err)
	}

	if err := s.addRoutes(cfg, reg); err != nil {
		return nil, fmt.Errorf("error adding routes: %w", err)
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

func (s *Server) addMiddlewares() error {
	s.engine.Use(func(c *gin.Context) {
		observer := metrics.ServerHttpRequestsDurationSeconds.WithLabelValues(c.Request.Method, s.srv.Addr, c.FullPath())

		timer := prometheus.NewTimer(observer)
		defer timer.ObserveDuration()

		metrics.ServerHttpRequestsTotal.WithLabelValues(c.Request.Method, s.srv.Addr, c.FullPath()).Inc()

		c.Next()
	})
	return nil
}

func (s *Server) addRoutes(cfg *config.Config, reg *registry.PluginRegistry) error {
	v1 := s.engine.Group("/v1")
	auth := s.engine.Group("/v1", tokenAuthMiddleware(cfg.Server.Token))

	v1.GET("/health", api.HandleHealth(reg))
	auth.GET("/plugins", api.HandlePlugins(reg))
	auth.GET("/services", api.HandleServices(reg))

	auth.GET("/queue", api.HandleGetQueue(cfg))
	auth.POST("/queue", api.HandlePostQueue(cfg))

	return nil
}

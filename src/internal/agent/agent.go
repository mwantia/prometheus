package agent

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/mwantia/prometheus/internal/configs"
	"github.com/mwantia/prometheus/internal/registry"
)

type PrometheusAgent struct {
	Mutex    sync.RWMutex
	Registry *registry.PluginRegistry
	Config   *configs.Config
}

func CreateNewAgent(c *configs.Config) *PrometheusAgent {
	return &PrometheusAgent{
		Registry: registry.NewRegistry(),
		Config:   c,
	}
}

func (a *PrometheusAgent) Serve(ctx context.Context) error {
	if err := a.serveLocalPlugins(); err != nil {
		log.Printf("Unable to serve local plugin: %v", err)
	}

	if err := a.serveEmbedPlugins(); err != nil {
		log.Printf("Unable to serve embed plugin: %v", err)
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	srv, err := a.startServer()
	if err != nil {
		return err
	}

	go a.Registry.Watch(ctx)

	go func() {
		log.Printf("Serving HTTP server: %s", a.Config.Agent.Server.Address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error serving http server: %v", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		log.Println("Shutting down agent...")

		shutdown := context.Background()
		shutdown, cncl := context.WithTimeout(shutdown, 10*time.Second)
		defer cncl()

		if err := srv.Shutdown(shutdown); err != nil {
			log.Fatalf("Error shutting down http server: %v", err)
		}
	}()

	wg.Wait()
	return a.Cleanup()
}

func (a *PrometheusAgent) Cleanup() error {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	var err error

	for _, plugin := range a.Registry.GetPlugins() {
		if shutdownErr := plugin.Plugin.Cleanup(); shutdownErr != nil {
			log.Printf("Error while performing cleanup for plugin '%s'", plugin.Name)
			err = errors.Join(err, shutdownErr)
		}

		if cleanupErr := plugin.Cleanup(); cleanupErr != nil {
			log.Printf("Error while performing cleanup for plugin '%s'", plugin.Name)
			err = errors.Join(err, cleanupErr)
		}
	}

	return err
}

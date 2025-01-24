package agent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/mwantia/prometheus/internal/agent/ops"
	"github.com/mwantia/prometheus/internal/config"
	"github.com/mwantia/prometheus/internal/registry"
)

type PrometheusAgent struct {
	Mutex    sync.RWMutex
	Registry *registry.PluginRegistry
	Config   *config.Config
}

func CreateNewAgent(c *config.Config) *PrometheusAgent {
	return &PrometheusAgent{
		Registry: registry.NewRegistry(),
		Config:   c,
	}
}

func (a *PrometheusAgent) Serve(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup
	cleanups := []func() error{}

	if err := a.serveLocalPlugins(); err != nil {
		log.Printf("Unable to serve local plugin: %v", err)
	}

	if err := a.serveEmbedPlugins(); err != nil {
		log.Printf("Unable to serve embed plugin: %v", err)
	}

	if a.Config.Server.Enabled {
		wg.Add(1)

		srv := ops.Server{}
		cleanup, err := srv.Create(a.Config, a.Registry)
		if err != nil {
			return fmt.Errorf("failed to create server: %w", err)
		}

		cleanups = append(cleanups, func() error {
			shutdown, cncl := context.WithTimeout(context.Background(), 10*time.Second)
			defer cncl()

			return cleanup(shutdown)
		})

		go func() {
			defer wg.Done()

			log.Println("Starting server...")
			if err := srv.Serve(ctx); err != nil {
				log.Fatalf("Error serving server: %v", err)
			}
		}()
	}

	if a.Config.Client.Enabled {
		wg.Add(1)

		client := ops.Client{}
		cleanup, err := client.Create(a.Config, a.Registry)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		cleanups = append(cleanups, func() error {
			shutdown, cncl := context.WithTimeout(context.Background(), 10*time.Second)
			defer cncl()

			return cleanup(shutdown)
		})

		go func() {
			defer wg.Done()

			log.Println("Starting client...")
			if err := client.Serve(ctx); err != nil {
				log.Fatalf("Error serving client: %v", err)
			}
		}()
	}

	if a.Config.Metrics.Enabled {
		wg.Add(1)

		metrics := ops.Metrics{}
		cleanup, err := metrics.Create(a.Config, a.Registry)
		if err != nil {
			return fmt.Errorf("failed to create metrics: %w", err)
		}

		cleanups = append(cleanups, func() error {
			shutdown, cncl := context.WithTimeout(context.Background(), 10*time.Second)
			defer cncl()

			return cleanup(shutdown)
		})

		go func() {
			defer wg.Done()

			log.Println("Starting metrics...")
			if err := metrics.Serve(ctx); err != nil {
				log.Fatalf("Error serving metrics: %v", err)
			}
		}()
	}

	go a.Registry.Watch(ctx)

	<-ctx.Done()
	log.Println("Shutting down agent...")

	for _, cleanup := range cleanups {
		if err := cleanup(); err != nil {
			log.Printf("Error during cleanup: %v", err)
		}
	}

	wg.Wait()
	return a.Cleanup()
}

func (a *PrometheusAgent) Cleanup() error {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	var err error

	for _, plugin := range a.Registry.GetPlugins() {
		if cleanupErr := plugin.Cleanup(); cleanupErr != nil {
			log.Printf("Error while performing cleanup for plugin '%s'", plugin.Name)
			err = errors.Join(err, cleanupErr)
		}
	}

	return err
}

package agent

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/mwantia/queueverse/internal/agent/ops"
	"github.com/mwantia/queueverse/internal/config"
	"github.com/mwantia/queueverse/internal/registry"
	"github.com/mwantia/queueverse/pkg/log"
)

type PrometheusAgent struct {
	Mutex    sync.RWMutex
	Log      log.Logger
	Registry *registry.PluginRegistry
	Config   *config.Config
}

func CreateNewAgent(c *config.Config) *PrometheusAgent {
	return &PrometheusAgent{
		Log:      log.New("agent"),
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
		a.Log.Warn("Unable to serve local plugin", "error", err)
	}

	if err := a.serveEmbedPlugins(); err != nil {
		a.Log.Warn("Unable to serve embed plugin", "error", err)
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

			a.Log.Debug("Starting server...")
			if err := srv.Serve(ctx); err != nil {
				a.Log.Error("Error serving server", "error", err)
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

			a.Log.Debug("Starting client...")
			if err := client.Serve(ctx); err != nil {
				a.Log.Error("Error serving client", "error", err)
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

			a.Log.Debug("Starting metrics...")
			if err := metrics.Serve(ctx); err != nil {
				a.Log.Error("Error serving metrics", "error", err)
			}
		}()
	}

	go a.Registry.Watch(ctx)

	<-ctx.Done()
	a.Log.Debug("Shutting down agent...")

	for _, cleanup := range cleanups {
		if err := cleanup(); err != nil {
			a.Log.Error("Error during agent cleanup", "error", err)
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
			a.Log.Error("Error while performing cleanup for plugin", "name", plugin.Name)
			err = errors.Join(err, cleanupErr)
		}
	}

	return err
}

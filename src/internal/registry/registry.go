package registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mwantia/prometheus/pkg/plugin"
)

type PluginRegistry struct {
	Mutex   sync.RWMutex
	Plugins map[string]*PluginInfo
}

type PluginInfo struct {
	Name           string                    `json:"name"`
	LastSeen       time.Time                 `json:"last_seen"`
	LastKnownError error                     `json:"-"`
	IsHealthy      bool                      `json:"is_healthy"`
	Capabilities   plugin.PluginCapabilities `json:"capabilities"`
	Plugin         plugin.Plugin             `json:"-"`
	Cleanup        PluginCleanup             `json:"-"`
}

type PluginCleanup func() error

func NewRegistry() *PluginRegistry {
	return &PluginRegistry{
		Plugins: make(map[string]*PluginInfo),
	}
}

func (r *PluginRegistry) Watch(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		plugins := r.GetPlugins()
		for _, plugin := range plugins {
			if err := plugin.Plugin.Health(); err != nil {
				plugin.IsHealthy = false
				plugin.LastKnownError = err

				continue
			}

			plugin.IsHealthy = true
			plugin.LastKnownError = nil
			plugin.LastSeen = time.Now()
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}

func (reg *PluginRegistry) RegisterPlugin(info *PluginInfo) error {
	reg.Mutex.Lock()
	defer reg.Mutex.Unlock()

	if _, exists := reg.Plugins[info.Name]; exists {
		return fmt.Errorf("a plugin with the name '%s' has already been registered", info.Name)
	}

	info.LastSeen = time.Now()
	reg.Plugins[info.Name] = info

	return nil
}

func (reg *PluginRegistry) Deregister(name string) (*PluginInfo, error) {
	reg.Mutex.Lock()
	defer reg.Mutex.Unlock()

	plugin, exists := reg.Plugins[name]
	if !exists {
		return nil, fmt.Errorf("a plugin with the name %s does not exist", name)
	}

	return plugin, nil
}

func (reg *PluginRegistry) GetPlugin(name string) (*PluginInfo, bool) {
	reg.Mutex.Lock()
	defer reg.Mutex.Unlock()

	plugin, exists := reg.Plugins[name]
	return plugin, exists
}

func (reg *PluginRegistry) GetPlugins() []*PluginInfo {
	reg.Mutex.Lock()
	defer reg.Mutex.Unlock()

	plugins := make([]*PluginInfo, 0, len(reg.Plugins))
	for _, plugin := range reg.Plugins {
		plugins = append(plugins, plugin)
	}

	return plugins
}

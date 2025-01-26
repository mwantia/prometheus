package registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mwantia/queueverse/pkg/plugin/identity"
	"github.com/mwantia/queueverse/pkg/plugin/tools"
)

type PluginRegistry struct {
	mutex   sync.RWMutex
	Plugins map[string]*PluginInfo
}

type PluginInfo struct {
	Name     string            `json:"name"`
	Version  string            `json:"version"`
	Author   string            `json:"author,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`

	LastSeen       time.Time `json:"last_seen"`
	LastKnownError error     `json:"-"`
	IsHealthy      bool      `json:"is_healthy"`

	Services PluginServices `json:"-"`
	Cleanup  PluginCleanup  `json:"-"`
}

type PluginServices struct {
	Identity identity.IdentityService `json:"-"`
	Tools    []tools.ToolService      `json:"-"`
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
			plugin.IsHealthy = true
			//
			for _, tool := range plugin.Services.Tools {
				if err := tool.Probe(); err != nil {
					name, _ := tool.GetName()

					plugin.IsHealthy = false
					plugin.LastKnownError = fmt.Errorf("error probing tool '%s': %w", name, err)

					continue
				}
			}

			if plugin.IsHealthy {
				plugin.LastKnownError = nil
				plugin.LastSeen = time.Now()
			}
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}

func (reg *PluginRegistry) RegisterPlugin(info *PluginInfo) error {
	reg.mutex.Lock()
	defer reg.mutex.Unlock()

	if _, exists := reg.Plugins[info.Name]; exists {
		return fmt.Errorf("a plugin with the name '%s' has already been registered", info.Name)
	}

	info.LastSeen = time.Now()
	reg.Plugins[info.Name] = info

	return nil
}

func (reg *PluginRegistry) Deregister(name string) (*PluginInfo, error) {
	reg.mutex.Lock()
	defer reg.mutex.Unlock()

	plugin, exists := reg.Plugins[name]
	if !exists {
		return nil, fmt.Errorf("a plugin with the name %s does not exist", name)
	}

	return plugin, nil
}

func (reg *PluginRegistry) GetPlugin(name string) (*PluginInfo, bool) {
	reg.mutex.Lock()
	defer reg.mutex.Unlock()

	plugin, exists := reg.Plugins[name]
	return plugin, exists
}

func (reg *PluginRegistry) GetPlugins() []*PluginInfo {
	reg.mutex.Lock()
	defer reg.mutex.Unlock()

	plugins := make([]*PluginInfo, 0, len(reg.Plugins))
	for _, plugin := range reg.Plugins {
		plugins = append(plugins, plugin)
	}

	return plugins
}

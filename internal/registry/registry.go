package registry

import (
	"fmt"

	"github.com/mwantia/queueverse/pkg/plugin/base"
)

func New() *Registry {
	return &Registry{
		Plugins: make(map[string]*RegistryPlugin),
	}
}

func (r *Registry) GetPlugins() []*RegistryPlugin {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	result := make([]*RegistryPlugin, 0, len(r.Plugins))
	for _, plugin := range r.Plugins {
		result = append(result, plugin)
	}

	return result
}

func (r *Registry) GetPluginInfo(name string) (*base.PluginInfo, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	plugin, exist := r.Plugins[name]
	if exist {
		return &plugin.Info, nil
	}

	return nil, fmt.Errorf("plugin with the name '%s' does not exist", name)
}

func (r *Registry) GetPluginStatus(name string) (*RegistryStatus, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	plugin, exist := r.Plugins[name]
	if exist {
		return &plugin.Status, nil
	}

	return nil, fmt.Errorf("plugin with the name '%s' does not exist", name)
}

func (r *Registry) Register(info *base.PluginInfo, impl interface{}, cleanup RegistryCleanup) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.Plugins[info.Name]; exists {
		return fmt.Errorf("a plugin with the name '%s' has already been registered", info.Name)
	}

	plugin := &RegistryPlugin{
		Info: *info,
		Status: RegistryStatus{
			IsHealthy: false,
		},
		Impl:    impl,
		Cleanup: cleanup,
	}

	r.Plugins[info.Name] = plugin
	return nil
}

func (r *Registry) Deregister(name string) (bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exist := r.Plugins[name]
	delete(r.Plugins, name)
	return exist, nil
}

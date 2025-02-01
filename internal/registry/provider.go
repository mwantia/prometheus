package registry

import (
	"fmt"

	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
)

func (r *Registry) GetProviders() ([]provider.ProviderPlugin, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	result := make([]provider.ProviderPlugin, 0)
	for _, plugin := range r.Plugins {
		if plugin.Info.Type == base.PluginProviderType {
			impl, success := plugin.Impl.(provider.ProviderPlugin)
			if !success {
				return nil, fmt.Errorf("failed to cast plugin '%s' as ProviderPlugin", plugin.Info.Name)
			}

			result = append(result, impl)
		}
	}

	return result, nil
}

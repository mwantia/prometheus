package registry

import (
	"fmt"
	"strings"

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

func (r *Registry) GetModelProvider(model string) (provider.ProviderPlugin, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, plugin := range r.Plugins {
		if plugin.Info.Type == base.PluginProviderType {
			// Try to cast impl as provider
			impl, success := plugin.Impl.(provider.ProviderPlugin)
			if !success {
				return nil, fmt.Errorf("failed to case plugin '%s' as 'provider.ProviderPlugin'", plugin.Info.Name)
			}

			models, err := impl.GetModels()
			if err != nil {
				return nil, fmt.Errorf("failed to get models for plugin '%s': %w", plugin.Info.Name, err)
			}

			for _, m := range *models {
				if strings.EqualFold(m.Name, model) {
					return impl, nil
				}
			}
		}
	}

	return nil, nil
}

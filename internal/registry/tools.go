package registry

import (
	"fmt"

	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/tools"
)

func (r *Registry) GetTools() ([]tools.ToolPlugin, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	result := make([]tools.ToolPlugin, 0)
	for _, plugin := range r.Plugins {
		if plugin.Type == base.PluginToolsType {
			impl, success := plugin.Impl.(tools.ToolPlugin)
			if !success {
				return nil, fmt.Errorf("failed to cast plugin '%s' as ToolPlugin", plugin.Info.Name)
			}

			result = append(result, impl)
		}
	}

	return result, nil
}

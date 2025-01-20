package plugin

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-plugin"
	"github.com/mwantia/prometheus/pkg/plugin/cache"
	"github.com/mwantia/prometheus/pkg/plugin/identity"
	"github.com/mwantia/prometheus/pkg/plugin/tools"
)

func ServeTools(is identity.IdentityService, services []tools.ToolService) error {
	plugins := map[string]plugin.Plugin{
		"identity": &identity.IdentityPlugin{
			Service: is,
		},
	}

	info, err := is.GetPluginInfo()
	if err != nil {
		return fmt.Errorf("error getting plugin info: %w", err)
	}

	for index, service := range services {
		name, _ := service.GetName()
		_, exist := identity.GetPluginServiceInfo(name, info.Services)
		if !exist {
			return fmt.Errorf("error receiving tool '%s'", name)
		}

		key := fmt.Sprintf("tool.%v", index)

		plugins[key] = &tools.ToolPlugin{
			Service: service,
		}
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins:         plugins,
	})
	return nil
}

func ServePlugin(i identity.IdentityService, c cache.CacheService) error {
	info, _ := i.GetPluginInfo()
	log.Printf("Serving plugin: %s", info.Name)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			"identity": &identity.IdentityPlugin{
				Service: i,
			},
			"cache": &cache.CachePlugin{
				Impl: c,
			},
		},
	})

	return nil
}

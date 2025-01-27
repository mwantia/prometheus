package identity

import "strings"

type IdentityService interface {
	GetPluginInfo() (*PluginInfo, error)
}

type PluginInfo struct {
	Name     string              `json:"name"`
	Version  string              `json:"version"`
	Author   string              `json:"author,omitempty"`
	Metadata map[string]string   `json:"metadata,omitempty"`
	Services []PluginServiceInfo `json:"services"`
}

type PluginServiceInfo struct {
	Name        string            `json:"service"`
	Type        PluginServiceType `json:"type"`
	Description string            `json:"description"`
}

type PluginServiceType string

const (
	IdentityServiceType PluginServiceType = "identity"
	ToolServiceType     PluginServiceType = "tool"
	CacheServiceType    PluginServiceType = "cache"
)

func GetPluginServiceInfo(name string, services []PluginServiceInfo) (*PluginServiceInfo, bool) {
	for _, service := range services {
		if strings.EqualFold(name, service.Name) {
			return &service, true
		}
	}

	return nil, false
}

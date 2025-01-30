package base

type BasePlugin interface {
	GetPluginInfo() (*PluginInfo, error)

	GetCapabilities() (*PluginCapabilities, error)

	SetConfig(*PluginConfig) error

	ProbePlugin() error
}

type PluginInfo struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Author  string `json:"author,omitempty"`
	Version string `json:"version,omitempty"`
}

type PluginConfig struct {
	ConfigMap map[string]interface{} `json:"-"`
}

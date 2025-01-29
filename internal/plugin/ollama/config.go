package ollama

import "github.com/mitchellh/mapstructure"

type OllamaPluginConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Model    string `mapstructure:"model"`
}

func (p *OllamaPlugin) setConfig(cfg map[string]interface{}) error {
	return mapstructure.Decode(cfg, &p.Config)
}

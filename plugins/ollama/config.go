package ollama

import "github.com/mitchellh/mapstructure"

type OllamaProviderConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Model    string `mapstructure:"model"`
}

func (p *OllamaProvider) setConfig(cfg map[string]interface{}) error {
	return mapstructure.Decode(cfg, &p.Config)
}

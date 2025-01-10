package ollama

import "github.com/mitchellh/mapstructure"

type OllamaConfig struct {
	Address string `mapstructure:"address,omitempty"`
}

func (p *OllamaPlugin) loadConfig(d map[string]interface{}) error {
	var cfg OllamaConfig
	if err := mapstructure.Decode(d, &cfg); err != nil {
		return err
	}

	p.Config = cfg
	return nil
}

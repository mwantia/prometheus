package anthropic

import "github.com/mitchellh/mapstructure"

type AnthropicConfig struct {
	Token string `mapstructure:"token"`
}

func (p *AnthropicProvider) setConfig(cfg map[string]interface{}) error {
	return mapstructure.Decode(cfg, &p.Config)
}

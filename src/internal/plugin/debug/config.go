package debug

import "github.com/mitchellh/mapstructure"

type DebugConfig struct {
	Foo string `mapstructure:"foo,omitempty"`
}

func (p *DebugPlugin) loadConfig(d map[string]interface{}) error {
	var cfg *DebugConfig
	if err := mapstructure.Decode(d, &cfg); err != nil {
		return err
	}

	p.Config = cfg
	return nil
}

package discord

import "github.com/mitchellh/mapstructure"

type DiscordConfig struct {
	AuthToken string `mapstructure:"auth_token,omitempty"`
}

func (p *DiscordPlugin) loadConfig(d map[string]interface{}) error {
	var cfg DiscordConfig
	if err := mapstructure.Decode(d, &cfg); err != nil {
		return err
	}

	p.Config = cfg
	return nil
}

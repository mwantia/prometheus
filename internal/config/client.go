package config

type ClientConfig struct {
	Enabled bool `hcl:"enabled,optional"`
}

func (c *ClientConfig) ValidateConfig() error {
	return nil
}

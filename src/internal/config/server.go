package config

import "fmt"

type ServerConfig struct {
	Enabled bool   `hcl:"enabled,optional"`
	Address string `hcl:"address,optional"`
}

func (c *ServerConfig) ValidateConfig() error {
	if c.Address == "" {
		return fmt.Errorf("'address' is required")
	}

	return nil
}

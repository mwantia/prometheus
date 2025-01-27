package config

import "fmt"

type MetricsConfig struct {
	Enabled bool   `hcl:"enabled,optional"`
	Address string `hcl:"address,optional"`
}

func (c *MetricsConfig) ValidateConfig() error {
	if c.Address == "" {
		return fmt.Errorf("'address' is required")
	}

	return nil
}

package config

import "fmt"

type RedisConfig struct {
	Address  string `hcl:"address,optional"`
	Database int    `hcl:"database,optional"`
}

func (c *RedisConfig) ValidateConfig() error {
	if c.Address == "" {
		return fmt.Errorf("address is required")
	}

	return nil
}

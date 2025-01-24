package config

import "fmt"

type RedisConfig struct {
	Endpoint string `hcl:"endpoint,optional"`
	Database int    `hcl:"database,optional"`
	Password string `hcl:"password,optional"`
}

func (c *RedisConfig) ValidateConfig() error {
	if c.Endpoint == "" {
		return fmt.Errorf("'endpoint' is required")
	}

	return nil
}

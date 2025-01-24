package config

import "fmt"

type ClientConfig struct {
	Enabled bool     `hcl:"enabled,optional"`
	Queues  []string `hcl:"queues,optional"`
}

func (c *ClientConfig) ValidateConfig() error {
	if c.Queues == nil || len(c.Queues) <= 0 {
		return fmt.Errorf("'queues' is required")
	}

	return nil
}

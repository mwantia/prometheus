package config

import "fmt"

type OllamaConfig struct {
	Address string `hcl:"address,optional"`
	Model   string `hcl:"model,optional"`
}

func (c *OllamaConfig) ValidateConfig() error {
	if c.Address == "" {
		return fmt.Errorf("address is required")
	}

	return nil
}

package config

import "fmt"

type OllamaConfig struct {
	Endpoint string `hcl:"endpoint,optional"`
	Model    string `hcl:"model,optional"`
}

func (c *OllamaConfig) ValidateConfig() error {
	if c.Endpoint == "" {
		return fmt.Errorf("'endpoint' is required")
	}

	return nil
}

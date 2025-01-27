package config

import "fmt"

type TelemetryConfig struct {
	Enabled     bool   `hcl:"enabled,optional"`
	Endpoint    string `hcl:"endpoint,optional"`
	ServiceName string `hcl:"servicename,optional"`
}

func (c *TelemetryConfig) ValidateConfig() error {
	if c.Endpoint == "" {
		return fmt.Errorf("'endpoint' is required")
	}
	if c.ServiceName == "" {
		return fmt.Errorf("'servicename' is required")
	}

	return nil
}

package config

import "fmt"

type Config struct {
	LogLevel     string   `hcl:"log_level,optional"`
	PluginDir    string   `hcl:"plugin_dir,optional"`
	EmbedPlugins []string `hcl:"embed_plugins,optional"`

	Server    *ServerConfig    `hcl:"server,block"`
	Client    *ClientConfig    `hcl:"client,block"`
	Metrics   *MetricsConfig   `hcl:"metrics,block"`
	Telemetry *TelemetryConfig `hcl:"telemetry,block"`
	Redis     *RedisConfig     `hcl:"redis,block"`
	Ollama    *OllamaConfig    `hcl:"ollama,block"`

	Plugins []*PluginConfig `hcl:"plugin,block"`
}

func CreateDefault() *Config {
	return &Config{
		LogLevel:     "INFO",
		PluginDir:    "./plugins",
		EmbedPlugins: make([]string, 0),

		Server: &ServerConfig{
			Enabled: true,
			Address: "0.0.0.0:8080",
		},
		Client: &ClientConfig{
			Enabled: true,
			Queues:  []string{"high", "normal", "low"},
		},
		Metrics: &MetricsConfig{
			Enabled: true,
			Address: "127.0.0.1:9001",
		},
		Telemetry: &TelemetryConfig{
			Enabled:     false,
			Endpoint:    "127.0.0.1:4318",
			ServiceName: "prometheus",
		},
		Redis: &RedisConfig{
			Endpoint: "127.0.0.1:6379",
			Database: 0,
			Password: "",
		},
		Ollama: &OllamaConfig{
			Endpoint: "127.0.0.1:11434",
		},

		Plugins: make([]*PluginConfig, 0),
	}
}

func (c *Config) ValidateConfig() error {
	if c.LogLevel == "" {
		return fmt.Errorf("'log_level' is required")
	}
	if c.PluginDir == "" {
		return fmt.Errorf("'plugin_dir' is required")
	}

	if c.Server == nil {
		return fmt.Errorf("block 'server' is required")
	}
	if err := c.Server.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid 'server' block: %w", err)
	}

	if c.Client == nil {
		return fmt.Errorf("block 'client' is required")
	}
	if err := c.Client.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid 'client' block: %w", err)
	}

	if c.Metrics == nil {
		return fmt.Errorf("block 'metrics' is required")
	}
	if err := c.Metrics.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid 'metrics' block: %w", err)
	}

	if c.Telemetry == nil {
		return fmt.Errorf("block 'telemetry' is required")
	}
	if err := c.Telemetry.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid 'telemetry' block: %w", err)
	}

	if c.Redis == nil {
		return fmt.Errorf("block 'redis' is required")
	}
	if err := c.Redis.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid 'redis' block: %w", err)
	}

	if c.Ollama == nil {
		return fmt.Errorf("block 'ollama' is required")
	}
	if err := c.Ollama.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid 'ollama' block: %w", err)
	}

	return nil
}

# QueueVerse

QueueVerse is a distributed queue system for Large Language Models (LLMs) with plugin support.

It provides a flexible architecture for managing LLM requests through a queue system while supporting various providers and tools through a plugin system.

## Features

- Distributed queue system for LLM requests
- Plugin-based architecture for extensibility
- Support for multiple LLM providers (Ollama, Anthropic)
- Built-in tool calling capabilities
- Metrics monitoring with Prometheus
- HTTP API for queue management
- Redis-based task persistence with `AsynQ`
- Health checking and monitoring

## Prerequisites

- Go 1.21 or higher
- Redis server
- Docker (for containerized deployment)

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/mwantia/queueverse.git

# Change to project directory
cd queueverse

# Download dependencies
go mod download

# Build the project
go build -o queueverse cmd/queueverse/main.go
```

### Using Docker

```bash
# Build the Docker image
docker build -t queueverse .

# Run the container
docker run -d \
  -p 8080:8080 \
  -p 9001:9001 \
  --name queueverse \
  queueverse
```

## Configuration

QueueVerse uses HCL format for configuration. Create a config file (e.g., `config.hcl`) with the following structure:

```hcl
log_level = "INFO"
pool_name = "default"
plugin_dir = "./plugins"
embed_plugins = ["ollama", "anthropic", "mock"]

server {
  enabled = true
  address = "0.0.0.0:8080"
  token = "your-auth-token"
}

client {
  enabled = true
}

metrics {
  enabled = true
  address = "127.0.0.1:9001"
}

redis {
  endpoint = "127.0.0.1:6379"
  database = 0
  password = ""
}

# Plugin configurations
plugin "ollama" {
  enabled = true
  config {
    endpoint = "http://localhost:11434"
  }
}

plugin "anthropic" {
  enabled = true
  config {
    token = "your-anthropic-api-token"
  }
}
```

## Usage

### Starting the Agent

```bash
# Start with default configuration
./queueverse agent

# Start with custom configuration
./queueverse agent --config /path/to/config.hcl
```

### API Endpoints

- `GET /v1/health` - Check service health
- `GET /v1/plugins` - List available plugins
- `GET /v1/models` - List available models
- `POST /v1/queue` - Submit a new task
- `GET /v1/queue/:task` - Get task status and result
- `GET /v1/queue` - List all tasks

### Example API Request

```bash
# Submit a new task
curl -X POST http://localhost:8080/v1/queue \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{
    "input": {
      "model": "claude-3-5-haiku-latest",
      "message": {
        "content": "Tell me the current time in Germany"
      },
      "style": "concise"
    }
  }'
```

## Plugin System

QueueVerse supports three types of plugins:

1. Provider Plugins: Interface with LLM providers (e.g., Ollama, Anthropic)
2. Tool Plugins: Implement specific functionalities (e.g., time queries, Discord integration)
3. Base Plugins: Core plugin functionality

### Built-in Plugins

- `mock`: Test plugin for development and testing
- `ollama`: Integration with Ollama LLM server
- `anthropic`: Integration with Anthropic's Claude models

### Creating Custom Plugins

To create a custom plugin:

1. Implement the required plugin interface (Provider, Tool, or Base)
2. Register the plugin in the configuration
3. Place the plugin binary in the configured plugin directory

Example plugin structure:

## Metrics

QueueVerse exposes Prometheus metrics at `:9001/metrics` including:

- HTTP request metrics
- Queue task metrics
- Plugin health metrics
- System metrics

## Development

### Project Structure

```yml
queueverse/
├── cmd/                  # Command line applications
├── internal/             # Internal packages
│   ├── agent/            # Agent implementation
│   ├── config/           # Configuration handling
│   ├── metrics/          # Metrics collection
│   └── registry/         # Plugin registry
├── pkg/                  # Public packages
│   ├── plugin/           # Plugin system
│   ├── tasks/            # Task definitions
│   └── log/              # Logging utilities
└── plugins/              # Built-in plugins
    ├── anthropic/        # Anthropic provider
    ├── ollama/           # Ollama provider
    └── mock/             # Mock provider
```

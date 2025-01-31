# QueueVerse

Ollama queue system and tool calling via plugin support.

## QueueVerse Plugins

```json
embed_plugins = ["<pluginname>"]
```

### Providers

Provider plugins are used to act as interface between `generate`-requests scheduled over **QueueVerse** and the endpoint responsible for generating the prompt result with the specified model.

#### Mock

Acts as a very simple dummy that can be accessed by specifying the model `mock-lorem-ipsum-8`, `mock-lorem-ipsum-16` or `mock-lorem-ipsum-32` which returns a fix number of sentences from **Lorem Ipsum**.

```json
plugin "mock" {
    config {

    }
}
```

#### Ollama

Connects to **Ollama** to perform prompt generation with the specified model defined in the request.

```json
plugin "ollama" {
    config {
        endpoint = "<ollama-http-endpoint>"
    }
}
```

#### Anthropic

```json
plugin "anthropic" {
    config {
        token = "<anthropic-api-token>
    }
}
```

### Tools

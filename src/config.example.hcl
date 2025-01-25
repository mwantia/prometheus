log_level     = "info"
pool_name     = "default"
plugin_dir    = "./plugins"
embed_plugins = []

server {
    enabled = true
    address = ":8080"
}

client {
    enabled = true
}

metrics {
    enabled = true
    address = ":9001"
}

redis {
    endpoint = "127.0.0.1:6379"
    database = 0
    password = ""
}

ollama {
    endpoint = "127.0.0.1:11434"
    model    = ""
}
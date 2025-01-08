log_level = "DEBUG"

agent {
    server {
        address = ":8080"
    }

    kafka {
        network   = "tcp"
        address   = "kafka.service.consul:9092"
        topics    = "discord-dm"
        partition = 0
    }
    
    plugin_dir    = "./tests/plugins"
    embed_plugins = ["debug"]
}

plugin "debug" {
    enabled = true

    config {
        foo = "bar"
    }
}
package plugin

import goplugin "github.com/hashicorp/go-plugin"

func ServePlugin(p Plugin) error {
	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]goplugin.Plugin{
			"driver": &PluginDriver{
				Impl: p,
			},
		},
	})

	return nil
}

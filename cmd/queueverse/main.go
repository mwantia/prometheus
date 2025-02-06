package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mwantia/queueverse/internal/agent"
	"github.com/mwantia/queueverse/internal/config"
	setuplog "github.com/mwantia/queueverse/internal/log"
	"github.com/spf13/cobra"
)

var (
	ConfigFlag  string
	NoColorFlag bool
)

var Config *config.Config

var (
	Root = &cobra.Command{
		Use:   "queueverse",
		Short: "LLM queue system and tool calling via plugin support",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.ParseConfig(ConfigFlag)
			if err != nil {
				return fmt.Errorf("unable to complete config: %w", err)
			}

			if err := cfg.ValidateConfig(); err != nil {
				return fmt.Errorf("unable to validate config: %w", err)
			}

			if err := setuplog.Setup(cfg.LogLevel); err != nil {
				return fmt.Errorf("unable to setup logging: %w", err)
			}

			if !strings.EqualFold("debug", cfg.LogLevel) {
				gin.SetMode(gin.ReleaseMode)
			}

			Config = cfg

			return nil
		},
	}
	Agent = &cobra.Command{
		Use:   "agent",
		Short: "Run QueueVerse agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return agent.CreateNewAgent(Config).Serve(context.Background())
		},
	}
	Plugin = &cobra.Command{
		Use:   "plugin [name]",
		Short: "Run embedded QueueVerse plugin",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("plugin name is required")
			}

			name := strings.TrimSpace(args[0])

			plugin, exists := agent.Plugins[name]
			if exists && plugin != nil {
				plugin()
				return nil
			}

			return fmt.Errorf("unknown plugin: %s", args[0])
		},
	}
)

func main() {
	Root.PersistentFlags().StringVar(&ConfigFlag, "config", "", "Defines the configuration path used by this application")
	Root.PersistentFlags().BoolVar(&NoColorFlag, "no-color", false, "Disables colored command output")

	Plugin.Flags().String("address", "", "If defined, registers the plugin in network mode and tries to connect to the external agent via 'address'.")

	Root.AddCommand(Agent, Plugin)

	if err := Root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

package discord

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/mwantia/prometheus/pkg/msg"
	"github.com/mwantia/prometheus/pkg/plugin"
)

type DiscordPlugin struct {
	plugin.DefaultPlugin

	Context context.Context
	Config  DiscordConfig
	Hub     msg.MessageHub

	Session *discordgo.Session
}

func NewPlugin() *DiscordPlugin {
	return &DiscordPlugin{
		Context: context.Background(),
	}
}

func (p *DiscordPlugin) Name() (string, error) {
	return "discord", nil
}

func (p *DiscordPlugin) Setup(s plugin.PluginSetup) error {
	p.Hub = s.Hub

	if err := p.loadConfig(s.Data); err != nil {
		log.Printf("Error converting mapstructure: %v", err)
	}

	gp, err := p.Hub.CreateProducer("global")
	if err != nil {
		return fmt.Errorf("error creating message hub producer")
	}

	if err := p.configureDiscordBot(gp); err != nil {
		return err
	}

	return nil
}

func (p *DiscordPlugin) Health() error {
	if p.Session == nil {
		return fmt.Errorf("discord session not initialized")
	}

	return nil
}

func (p *DiscordPlugin) Cleanup() error {
	if err := p.Session.Close(); err != nil {
		return fmt.Errorf("unable to close discord session: %v", err)
	}

	return nil
}

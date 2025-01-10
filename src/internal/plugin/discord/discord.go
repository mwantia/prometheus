package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/mwantia/prometheus/pkg/msg"
)

func (p *DiscordPlugin) configureDiscordBot(gp msg.MessageHubProducer) error {
	session, err := discordgo.New("Bot " + p.Config.AuthToken)
	if err != nil {
		return err
	}

	session.AddHandler(p.handleMessageCreate(gp))
	session.Identify.Intents = discordgo.IntentDirectMessages

	p.Session = session
	return session.Open()
}

func (p *DiscordPlugin) handleMessageCreate(gp msg.MessageHubProducer) interface{} {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			log.Printf("Error getting channel: %v", err)
			return
		}

		if err := s.ChannelTyping(channel.ID); err != nil {
			log.Printf("Error showing typing indicator: %v", err)
		}

		if channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM {
			gc, err := p.Hub.CreateConsumer(m.ChannelID)
			if err != nil {
				log.Printf("Error creating message hub consumer")
			}
			defer gc.Cleanup(p.Context)

			if err := gp.Write(p.Context, m.Content); err != nil {
				log.Printf("Unable to write discord content to message hub: %v", err)
			}

			if err := gc.Read(p.Context, p.handleMessageRespond(s, channel)); err != nil {
				log.Printf("Error handling consumer messages: %v", err)
			}

			return
		}
	}
}

func (p *DiscordPlugin) handleMessageRespond(s *discordgo.Session, c *discordgo.Channel) interface{} {
	return func(r string) {
		if _, err := s.ChannelMessageSend(c.ID, r); err != nil {
			log.Printf("Error sending message to channel '%s': %v", c.ID, err)
		}

		if r != "" {
			if err := s.ChannelTyping(c.ID); err != nil {
				log.Printf("Error showing typing indicator: %v", err)
			}
		}
	}
}

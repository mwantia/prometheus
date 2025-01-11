package discord

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mwantia/prometheus/pkg/msg"
)

func (p *DiscordPlugin) configureDiscordBot() error {
	session, err := discordgo.New("Bot " + p.Config.AuthToken)
	if err != nil {
		return err
	}

	session.AddHandler(p.handleMessageCreate())
	session.Identify.Intents = discordgo.IntentDirectMessages

	p.Session = session
	return session.Open()
}

func (p *DiscordPlugin) handleMessageCreate() interface{} {
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
			consumer, err := p.Hub.CreateConsumer("conversations." + m.ChannelID)
			if err != nil {
				log.Printf("Unable to create 'conversations.<ChannelID>': %v", err)
			}

			if err := p.sendConversationCreateMessage(m.ID, channel); err != nil {
				log.Printf("Unable to write to 'conversations.create': %v", err)
			}

			if err := consumer.Read(p.Context, p.handleMessageRespond(s, channel)); err != nil {
				log.Printf("Unable to read 'conversations.<ChannelID>': %v", err)
			}
		}
	}
}

func (p *DiscordPlugin) sendConversationCreateMessage(id string, c *discordgo.Channel) error {
	producer, err := p.Hub.CreateProducer("conversations.create")
	if err != nil {
		return err
	}

	return producer.Write(p.Context, msg.Message{
		ID:        id,
		Content:   c.ID,
		Type:      "msg",
		Sequence:  0,
		Timestamp: time.Now(),
	})
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

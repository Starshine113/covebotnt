package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Ping command: replies with latency and message edit time
func commandPing(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	heartbeat := s.HeartbeatLatency().String()

	// get current time
	cmdStart := time.Now()

	// send initial message
	message, err := s.ChannelMessageSend(m.ChannelID, "Pong!\nHeartbeat: "+heartbeat)
	if err != nil {
		return fmt.Errorf("Ping: %w", err)
	}

	// get time difference, edit message
	diff := time.Now().Sub(cmdStart).String()
	_, err = s.ChannelMessageEdit(message.ChannelID, message.ID, message.Content+"\nMessage latency: "+diff)
	return err
}

// Help shows the help pages
func commandHelp(args []string, s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	if len(args) == 0 {
		embed := &discordgo.MessageEmbed{
			Title:       "CoveBotn't help",
			Description: "CoveBotn't is a general purpose bot, with a gatekeeper, moderation commands, and starboard functionality.",
			Color:       0x21a1a8,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Created by Starshine System (Starshine ☀✨#5000) | CoveBotn't v0.1",
			},
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Basic commands", Value: "`ping`: show the bot's latency\n`help`: show this help page", Inline: false},
			},
		}

		_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			return fmt.Errorf("Help: %w", err)
		}
	}
	return nil
}

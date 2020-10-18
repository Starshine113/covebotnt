package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func commandStarboard(args []string, s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	// this command needs the mod role or administrator perms
	perms := checkModRole(s, m.Author.ID, m.GuildID, false)
	if perms != nil {
		commandError(perms, s, m)
		return
	}
	if len(args) == 0 {
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, currentStarboardSettings(m.GuildID))
		if err != nil {
			return fmt.Errorf("Starboard: %w", err)
		}
	}
	_, err = s.ChannelMessageSend(m.ChannelID, "congrats you have permission to do this")
	if err != nil {
		return fmt.Errorf("Starboard: %w", err)
	}
	return
}

func currentStarboardSettings(guildID string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{}
}

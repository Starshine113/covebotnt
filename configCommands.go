package main

import "github.com/bwmarrin/discordgo"

func commandStarboard(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	// this command needs the mod role or administrator perms
	perms := checkModRole(s, m.Author.ID, m.GuildID, false)
	if perms != nil {
		commandError(perms, s, m)
		return
	}
	if len(args) == 0 {
		s.ChannelMessageSendEmbed(m.ChannelID, currentStarboardSettings(m.GuildID))
	}
	s.ChannelMessageSend(m.ChannelID, "congrats you have permission to do this")
}

func currentStarboardSettings(guildID string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{}
}

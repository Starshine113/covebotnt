package main

import "github.com/bwmarrin/discordgo"

func commandError(err error, s *discordgo.Session, m *discordgo.MessageCreate) {
	s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
	s.ChannelMessageSend(m.ChannelID, "❌ An error occured:\n> "+err.Error())
	sugar.Errorf("Command error occured: ", err.Error())
}

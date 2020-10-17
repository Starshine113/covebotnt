package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func commandSetStatus(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	// this command needs bot owner permissions
	perms := checkOwner(m.Author.ID)
	if perms != nil {
		commandError(perms, s, m)
		return
	}

	// this command needs at least 1 argument
	if len(args) == 0 {
		commandError(&errorNotEnoughArgs{1, 0}, s, m)
		return
	}

	// set the status to the specified string
	status := strings.Join(args, " ")
	dg.UpdateStatus(0, status)
	s.ChannelMessageSend(m.ChannelID, "Set status to `+"+status+"`")
}

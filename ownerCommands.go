package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func commandSetStatus(args []string, s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
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
	_, err = s.ChannelMessageSend(m.ChannelID, "Set status to `+"+status+"`")
	if err != nil {
		return fmt.Errorf("SetStatus: %v", err)
	}
	return
}

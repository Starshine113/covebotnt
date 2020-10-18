package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func commandGetUser(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	perms := checkOwner(m.Author.ID)
	if perms != nil {
		commandError(perms, s, m)
		return
	}
	if len(args) != 1 {
		commandError(&errorNotEnoughArgs{2, len(args)}, s, m)
		return
	}
	user, err := parseUser(s, args[0])
	if err != nil {
		commandError(err, s, m)
	}
	s.ChannelMessageSend(m.ChannelID, user.Username)
}

func commandSetStatus(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	// this command needs bot owner permissions
	perms := checkOwner(m.Author.ID)
	if perms != nil {
		commandError(perms, s, m)
		return
	}

	// this command needs at least 2 arguments
	if len(args) < 2 {
		commandError(&errorNotEnoughArgs{2, len(args)}, s, m)
		return
	}

	// check the first arg -- boolean between -replace and -append
	if args[0] == "-replace" {
		config.Bot.CustomStatus.Override = true
	} else if args[0] == "-append" {
		config.Bot.CustomStatus.Override = false
	} else {
		commandError(&errorMissingRequiredArgs{"<-replace/-append> [-clear] <status string>", "<-replace/-append>"}, s, m)
		return
	}

	// set custom status to the specified string
	config.Bot.CustomStatus.Status = strings.Join(args[1:], " ")

	// check second argument -- if it's `-clear` the custom status is cleared
	if args[1] == "-clear" {
		config.Bot.CustomStatus.Status = ""
	}

	// set the status
	newStatus := config.Bot.Prefixes[0] + "help | in " + fmt.Sprint(len(s.State.Guilds)) + " guilds"
	if config.Bot.CustomStatus.Status != "" {
		newStatus += " | " + config.Bot.CustomStatus.Status
	}
	if config.Bot.CustomStatus.Override {
		newStatus = config.Bot.CustomStatus.Status
	}
	err := dg.UpdateStatus(0, newStatus)
	if err != nil {
		commandError(err, s, m)
		return
	}
	_, err = s.ChannelMessageSend(m.ChannelID, "Set status to `"+newStatus+"`")
	if err != nil {
		sugar.Errorf("Error when sending message: ", err)
		return
	}
	sugar.Infof("Updated status to \"%v\"", newStatus)
}

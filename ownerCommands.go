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
		return nil
	}

	// this command needs at least 2 arguments
	if len(args) < 2 {
		commandError(&errorNotEnoughArgs{2, len(args)}, s, m)
		return nil
	}

	// check the first arg -- boolean between -replace and -append
	if args[0] == "-replace" {
		config.Bot.CustomStatus.Override = true
	} else if args[0] == "-append" {
		config.Bot.CustomStatus.Override = false
	} else {
		commandError(&errorMissingRequiredArgs{"<-replace/-append> [-clear] <status string>", "<-replace/-append>"}, s, m)
		return nil
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
	err = dg.UpdateStatus(0, newStatus)
	if err != nil {
		commandError(err, s, m)
		return err
	}
	_, err = s.ChannelMessageSend(m.ChannelID, "Set status to `"+newStatus+"`")
	if err != nil {
		sugar.Errorf("Error when sending message: ", err)
		return err
	}
	sugar.Infof("Updated status to \"%v\"", newStatus)
	return nil
}

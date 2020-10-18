package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func commandEcho(args []string, s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	perms := checkModRole(s, m.Author.ID, m.GuildID, false)
	if perms != nil {
		commandError(perms, s, m)
		return nil
	}

	if len(args) == 0 {
		commandError(&errorNotEnoughArgs{1, len(args)}, s, m)
		return nil
	}

	channelID := m.ChannelID

	// check if the user supplied a -channel argument
	if args[0] == "-channel" || args[0] == "-chan" || args[0] == "-ch" {
		if len(args) == 1 {
			commandError(&errorNotEnoughArgs{2, len(args)}, s, m)
			return nil
		}
		channel, err := parseChannel(s, args[1], m.GuildID)
		if err != nil {
			commandError(err, s, m)
			return nil
		}
		channelID = channel.ID
		args = args[2:]
	}

	message := strings.Join(args, " ")
	_, err = s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: message,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
		},
	})
	if err != nil {
		return err
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, successEmoji)
	if err != nil {
		return nil
	}

	if channelID != m.ChannelID {
		_, err = s.ChannelMessageSend(m.ChannelID, successEmoji+" Sent message to <#"+channelID+">")
		if err != nil {
			return err
		}
	}
	return nil
}

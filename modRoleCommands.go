package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func commandPrefix(args []string, s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	perms := checkModRole(s, m.Author.ID, m.GuildID, false)
	if perms != nil {
		commandError(perms, s, m)
		return nil
	}

	// if there are no arguments, show the current prefix
	if len(args) == 0 {
		if globalSettings[m.GuildID].Prefix == "" {
			_, err = s.ChannelMessageSend(m.ChannelID, "The current prefix is `"+config.Bot.Prefixes[0]+"` (default).")
			if err != nil {
				return err
			}
		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, "The current prefix is `"+globalSettings[m.GuildID].Prefix+"`.")
			if err != nil {
				return err
			}
		}
		return nil
	}

	// if there's more than 1 argument, error
	if len(args) > 1 {
		commandError(&errorTooManyArguments{1, len(args)}, s, m)
		return nil
	}

	// otherwise, set prefix to first argument
	err = setGuildPrefix(args[0], m.GuildID)
	if err != nil {
		commandError(err, s, m)
		return nil
	}
	_, err = s.ChannelMessageSend(m.ChannelID, "Changed prefix to `"+globalSettings[m.GuildID].Prefix+"`.")
	if err != nil {
		return err
	}

	return nil
}

func commandModRoles(args []string, s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	perms := checkAdmin(s, m.Author.ID, m.GuildID)
	if perms != nil {
		commandError(perms, s, m)
		return nil
	}

	if len(args) >= 1 {
		switch args[0] {
		case "add":
			_, err := s.ChannelMessageSend(m.ChannelID, "add a role")
			if err != nil {
				commandError(err, s, m)
				return nil
			}
		case "remove", "delete":
			_, err := s.ChannelMessageSend(m.ChannelID, "remove a role")
			if err != nil {
				commandError(err, s, m)
				return nil
			}
		case "list", "get":
			_, err := s.ChannelMessageSend(m.ChannelID, "list roles")
			if err != nil {
				commandError(err, s, m)
				return nil
			}
		}
	} else {
		roleEmbed, err := getModRoles(s, m.GuildID)
		if err != nil {
			commandError(err, s, m)
			return nil
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, roleEmbed)
		if err != nil {
			return fmt.Errorf("ModRoles: %v", err)
		}
	}
	return nil
}

func getModRoles(s *discordgo.Session, guildID string) (*discordgo.MessageEmbed, error) {
	var roles []string
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return nil, err
	}
	for _, modRole := range globalSettings[guildID].Moderation.ModRoles {
		roles = append(roles, "<@&"+modRole+">")
	}

	var roleString string
	if len(roles) == 0 {
		roleString = "No mod roles."
	} else {
		for i, role := range roles {
			if i == len(roles)-1 {
				roleString += role
			} else {
				roleString += role + ", "
			}
		}
	}
	embed := &discordgo.MessageEmbed{
		Title:       "Mod roles for " + guild.Name,
		Description: roleString,
		Color:       0x21a1a8,
	}
	return embed, nil
}

package main

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func commandStarboard(args []string, s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	// this command needs the mod role or administrator perms
	perms := checkModRole(s, m.Author.ID, m.GuildID, false)
	if perms != nil {
		commandError(perms, s, m)
		return nil
	}

	guild, err := s.Guild(m.GuildID)
	if err != nil {
		commandError(err, s, m)
		return nil
	}

	if len(args) == 0 {
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, currentStarboardSettings(guild))
		if err != nil {
			return fmt.Errorf("Starboard: %w", err)
		}
	}
	if len(args) == 2 {
		if args[0] == "channel" {
			channel, err := parseChannel(s, args[1], m.GuildID)
			if err != nil {
				commandError(err, s, m)
				return nil
			}
			err = setStarboardChannel(channel.ID, m.GuildID)
			if err != nil {
				commandError(err, s, m)
				return nil
			}
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Title:       "Starboard channel changed",
				Description: "Changed the starboard channel for " + guild.Name + " to " + channel.Mention(),
			})
			if err != nil {
				return err
			}
		} else if args[0] == "limit" {
			limit, err := strconv.ParseInt(args[1], 10, 0)
			if err != nil {
				commandError(err, s, m)
				return nil
			}
			err = setStarboardLimit(int(limit), m.GuildID)
			if err != nil {
				commandError(err, s, m)
				return nil
			}
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Title:       "Starboard limit changed",
				Description: "Changed the starboard limit for " + guild.Name + " to " + fmt.Sprint(limit),
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func currentStarboardSettings(guild *discordgo.Guild) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Current starboard settings for " + guild.Name,
		Description: "Starboard channel is <#" + globalSettings[guild.ID].Starboard.StarboardID + ">\nThe star emoji is " + globalSettings[guild.ID].Starboard.Emoji + "\nThe current requirement is " + fmt.Sprint(globalSettings[guild.ID].Starboard.ReactLimit) + " reactions",
	}
}

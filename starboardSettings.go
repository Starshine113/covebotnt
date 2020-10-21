package main

import (
	"fmt"
	"strconv"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/bwmarrin/discordgo"
)

func commandStarboard(ctx *cbctx.Ctx) (err error) {
	err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
	if err != nil {
		return err
	}

	// this command needs the mod role or administrator perms
	perms := checkModRole(ctx.Session, ctx.Author.ID, ctx.Message.GuildID, false)
	if perms != nil {
		commandError(perms, ctx.Session, ctx.Message)
		return nil
	}

	guild, err := ctx.Session.Guild(ctx.Message.GuildID)
	if err != nil {
		commandError(err, ctx.Session, ctx.Message)
		return nil
	}

	if len(ctx.Args) == 0 {
		_, err = ctx.Send(currentStarboardSettings(guild))
		if err != nil {
			return fmt.Errorf("Starboard: %w", err)
		}
	}
	if len(ctx.Args) == 2 {
		if ctx.Args[0] == "channel" {
			channel, err := ctx.ParseChannel(ctx.Args[1])
			if err != nil {
				commandError(err, ctx.Session, ctx.Message)
				return nil
			}
			err = setStarboardChannel(channel.ID, ctx.Message.GuildID)
			if err != nil {
				commandError(err, ctx.Session, ctx.Message)
				return nil
			}
			_, err = ctx.Send(&discordgo.MessageEmbed{
				Title:       "Starboard channel changed",
				Description: "Changed the starboard channel for " + guild.Name + " to " + channel.Mention(),
			})
			if err != nil {
				return err
			}
		} else if ctx.Args[0] == "limit" {
			limit, err := strconv.ParseInt(ctx.Args[1], 10, 0)
			if err != nil {
				commandError(err, ctx.Session, ctx.Message)
				return nil
			}
			err = setStarboardLimit(int(limit), ctx.Message.GuildID)
			if err != nil {
				commandError(err, ctx.Session, ctx.Message)
				return nil
			}
			_, err = ctx.Send(&discordgo.MessageEmbed{
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

package main

import (
	"fmt"
	"strconv"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/Starshine113/covebotnt/structs"
	"github.com/bwmarrin/discordgo"
)

func commandStarboard(ctx *crouter.Ctx) (err error) {
	err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
	if err != nil {
		return err
	}

	guild, err := ctx.Session.Guild(ctx.Message.GuildID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	s, err := ctx.Database.GetGuildSettings(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	_, err = ctx.Send(currentStarboardSettings(guild, s))
	if err != nil {
		return fmt.Errorf("Starboard: %w", err)
	}
	return nil
}

func commandStarboardChannel(ctx *crouter.Ctx) (err error) {
	if len(ctx.Args) == 0 {
		ctx.CommandError(&crouter.ErrorNotEnoughArgs{
			NumRequiredArgs: 1,
			SuppliedArgs:    0,
		})
		return nil
	}
	guild, err := ctx.Session.Guild(ctx.Message.GuildID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	channel, err := ctx.ParseChannel(ctx.Args[0])
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	err = pool.SetStarboardChannel(channel.ID, ctx.Message.GuildID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	_, err = ctx.Send(&discordgo.MessageEmbed{
		Title:       "Starboard channel changed",
		Description: "Changed the starboard channel for " + guild.Name + " to " + channel.Mention(),
	})
	return
}

func commandStarboardLimit(ctx *crouter.Ctx) (err error) {
	if len(ctx.Args) == 0 {
		ctx.CommandError(&crouter.ErrorNotEnoughArgs{
			NumRequiredArgs: 1,
			SuppliedArgs:    0,
		})
		return nil
	}
	guild, err := ctx.Session.Guild(ctx.Message.GuildID)
	limit, err := strconv.ParseInt(ctx.Args[0], 10, 0)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	err = pool.SetStarboardLimit(int(limit), ctx.Message.GuildID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	_, err = ctx.Send(&discordgo.MessageEmbed{
		Title:       "Starboard limit changed",
		Description: "Changed the starboard limit for " + guild.Name + " to " + fmt.Sprint(limit),
	})
	return
}

func commandStarboardEmoji(ctx *crouter.Ctx) (err error) {
	return
}

func currentStarboardSettings(guild *discordgo.Guild, settings structs.GuildSettings) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Current starboard settings for " + guild.Name,
		Description: "Starboard channel is <#" + settings.Starboard.StarboardID + ">\nThe star emoji is " + settings.Starboard.Emoji + "\nThe current requirement is " + fmt.Sprint(settings.Starboard.ReactLimit) + " reactions",
	}
}

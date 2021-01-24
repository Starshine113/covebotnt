package main

import (
	"github.com/starshine-sys/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

func commandModLogChannel(ctx *crouter.Ctx) (err error) {
	err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
	if err != nil {
		return err
	}

	if len(ctx.Args) != 1 {
		ctx.CommandError(&crouter.ErrorMissingRequiredArgs{RequiredArgs: "channel", MissingArgs: "channel"})
		return nil
	}

	channel, err := ctx.ParseChannel(ctx.Args[0])
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	err = ctx.Database.SetModLogChannel(ctx.Message.GuildID, channel.ID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	err = ctx.Database.RemoveFromGuildCache(ctx.Message.GuildID)
	if err != nil {
		return err
	}

	_, err = ctx.Send(&discordgo.MessageEmbed{
		Title:       "Mod logs channel changed",
		Description: "Changed the mod logs channel to " + channel.Mention(),
	})
	return
}

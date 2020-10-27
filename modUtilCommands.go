package main

import (
	"strings"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/Starshine113/flagparser"
	"github.com/bwmarrin/discordgo"
)

func commandEcho(ctx *cbctx.Ctx) (err error) {
	err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
	if err != nil {
		return err
	}

	if len(ctx.Args) == 0 {
		ctx.CommandError(&errorNotEnoughArgs{1, len(ctx.Args)})
		return nil
	}

	channelID := ctx.Message.ChannelID

	flagParser, err := flagparser.NewFlagParser(flagparser.String("channel", "chan", "ch"))
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	processedArgs, err := flagParser.Parse(ctx.Args)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	ctx.Args = processedArgs["rest"].([]string)

	// check if the user supplied a -channel argument
	if processedArgs["channel"].(string) != "" {
		channel, err := ctx.ParseChannel(processedArgs["channel"].(string))
		if err != nil {
			ctx.CommandError(err)
			return nil
		}
		channelID = channel.ID
	}

	err = ctx.Session.ChannelTyping(channelID)
	if err != nil {
		return err
	}

	message := strings.Join(ctx.Args, " ")
	_, err = ctx.Session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: message,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
		},
	})
	if err != nil {
		return err
	}

	err = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, cbctx.SuccessEmoji)
	if err != nil {
		return nil
	}

	if channelID != ctx.Message.ChannelID {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, cbctx.SuccessEmoji+" Sent message to <#"+channelID+">")
		if err != nil {
			return err
		}
	}
	return nil
}

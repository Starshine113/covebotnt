package commands

import (
	"errors"
	"strings"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/Starshine113/flagparser"
	"github.com/bwmarrin/discordgo"
)

// Echo says whatever the user inputs through the bot
func Echo(ctx *crouter.Ctx) (err error) {
	if len(ctx.Args) == 0 {
		ctx.CommandError(&crouter.ErrorNotEnoughArgs{NumRequiredArgs: 1, SuppliedArgs: len(ctx.Args)})
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
		if channel.GuildID != ctx.Message.GuildID {
			ctx.CommandError(errors.New("you cannot echo messages into a channel in a different server"))
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

	err = ctx.React(crouter.SuccessEmoji)
	if err != nil {
		return nil
	}

	if channelID != ctx.Message.ChannelID {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, crouter.SuccessEmoji+" Sent message to <#"+channelID+">")
		if err != nil {
			return err
		}
	}
	return nil
}

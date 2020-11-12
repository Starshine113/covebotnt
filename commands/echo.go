package commands

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/Starshine113/flagparser"
	"github.com/bwmarrin/discordgo"
)

// Echo says whatever the user inputs through the bot
func Echo(ctx *crouter.Ctx) (err error) {
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

	if len(ctx.Args) == 0 && len(ctx.Message.Attachments) == 0 {
		_, err = ctx.CommandError(&crouter.ErrorNotEnoughArgs{
			NumRequiredArgs: 1,
			SuppliedArgs:    0})
		return err
	}

	// check if the user supplied a -channel argument
	if processedArgs["channel"].(string) != "" {
		channel, err := ctx.ParseChannel(processedArgs["channel"].(string))
		if err != nil {
			_, err = ctx.CommandError(err)
			return nil
		}
		if channel.GuildID != ctx.Message.GuildID {
			_, err = ctx.CommandError(errors.New("you cannot echo messages into a channel in a different server"))
			return nil
		}
		channelID = channel.ID
	}

	err = ctx.Session.ChannelTyping(channelID)
	if err != nil {
		return err
	}

	var attachments []*discordgo.File
	if len(ctx.Message.Attachments) != 0 {
		for _, a := range ctx.Message.Attachments {
			resp, err := http.Get(a.URL)
			if err != nil {
				_, err = ctx.CommandError(err)
				return err
			}
			defer resp.Body.Close()

			file := discordgo.File{
				Name:   a.Filename,
				Reader: resp.Body,
			}

			attachments = append(attachments, &file)
		}
	}

	message := strings.Join(ctx.Args, " ")
	_, err = ctx.Session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: message,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
		},
		Files: attachments,
	})
	if err != nil {
		return err
	}

	ctx.React(crouter.SuccessEmoji)
	if channelID != ctx.Message.ChannelID {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, crouter.SuccessEmoji+" Sent message to <#"+channelID+">")
		if err != nil {
			return err
		}
	}
	return nil
}

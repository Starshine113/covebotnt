package modutilcommands

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
	flag "github.com/spf13/pflag"
)

// Echo says whatever the user inputs through the bot
func Echo(ctx *crouter.Ctx) (err error) {
	fs := flag.NewFlagSet("flags", flag.ContinueOnError)
	delete := fs.BoolP("delete", "d", false, "Whether or not to automatically delete the command")
	out := fs.StringP("channel", "c", "", "The channel to send the output to")

	err = fs.Parse(ctx.Args)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	ctx.Args = fs.Args()

	if len(ctx.Args) == 0 && len(ctx.Message.Attachments) == 0 {
		_, err = ctx.CommandError(&crouter.ErrorNotEnoughArgs{
			NumRequiredArgs: 1,
			SuppliedArgs:    0})
		return err
	}

	if *delete {
		ctx.Session.ChannelMessageDelete(ctx.Channel.ID, ctx.Message.ID)
	}

	// check if the user supplied a channel
	if *out != "" {
		channel, err := ctx.ParseChannel(*out)
		if err != nil {
			_, err = ctx.CommandError(err)
			return nil
		}
		if channel.GuildID != ctx.Message.GuildID {
			_, err = ctx.CommandError(errors.New("you cannot echo messages into a channel in a different server"))
			return nil
		}
		*out = channel.ID
	} else {
		*out = ctx.Channel.ID
	}

	err = ctx.Session.ChannelTyping(*out)
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

	var m *discordgo.MessageAllowedMentions
	if ctx.Message.MentionEveryone {
		m = &discordgo.MessageAllowedMentions{Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers, discordgo.AllowedMentionTypeEveryone, discordgo.AllowedMentionTypeRoles}}
	} else {
		m = &discordgo.MessageAllowedMentions{Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers}}
	}

	message := strings.Join(ctx.Args, " ")

	_, err = ctx.Session.ChannelMessageSendComplex(*out, &discordgo.MessageSend{
		Content:         message,
		AllowedMentions: m,
		Files:           attachments,
	})
	if err != nil {
		return err
	}

	ctx.React(crouter.SuccessEmoji)
	if *out != ctx.Message.ChannelID {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, crouter.SuccessEmoji+" Sent message to <#"+*out+">")
		if err != nil {
			return err
		}
	}
	return nil
}

package ownercommands

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// AdminDM sends a message to a user through the bot
func AdminDM(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(2); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	_, err = discordgo.SnowflakeTimestamp(ctx.Args[0])
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	msg := strings.Join(ctx.Args[1:], " ")

	var attachmentLink string
	if len(ctx.Message.Attachments) > 0 {
		match, _ := regexp.MatchString("\\.(png|jpg|jpeg|gif|webp)$", ctx.Message.Attachments[0].URL)
		if match {
			attachmentLink = ctx.Message.Attachments[0].URL
		}
	}

	embed := &discordgo.MessageEmbed{
		Description: msg,
		Color:       0x21a1a8,
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: ctx.BotUser.AvatarURL("128"),
			Name:    ctx.BotUser.Username + " Admin",
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Message ID: " + ctx.Message.ID,
		},
		Image: &discordgo.MessageEmbedImage{
			URL: attachmentLink,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	ch, err := ctx.Session.UserChannelCreate(ctx.Args[0])
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	send := &discordgo.MessageSend{
		Embed: embed,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{},
		},
	}

	_, err = ctx.Session.ChannelMessageSendComplex(ch.ID, send)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	send.Content = fmt.Sprintf("> Successfully sent message to user %v", ctx.Args[0])
	_, err = ctx.Send(send)
	if err != nil {
		return err
	}

	err = ctx.React(crouter.SuccessEmoji)
	return
}

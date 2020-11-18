package commands

import (
	"fmt"
	"strings"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// DM sends a message *without hiding the sender* through the bot
func DM(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(2); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	m, err := ctx.ParseMember(ctx.Args[0])
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	guild, err := ctx.Session.Guild(ctx.Message.GuildID)

	msg := strings.Join(ctx.Args[1:], " ")

	embed := &discordgo.MessageEmbed{
		Description: msg,
		Color:       0x21a1a8,
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: ctx.Author.AvatarURL("256"),
			Name:    ctx.Author.String(),
		},
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: guild.IconURL(),
			Text:    "Message ID: " + ctx.Message.ID,
		},
	}

	ch, err := ctx.Session.UserChannelCreate(m.User.ID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	send := &discordgo.MessageSend{
		Content: fmt.Sprintf("You have received a message from a moderator in **%v**:", guild.Name),
		Embed:   embed,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{},
		},
	}

	_, err = ctx.Session.ChannelMessageSendComplex(ch.ID, send)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	err = ctx.React(crouter.SuccessEmoji)
	return
}

// AnonDM sends a message, *hiding the sender*, through the bot
func AnonDM(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(2); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	m, err := ctx.ParseMember(ctx.Args[0])
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	guild, err := ctx.Session.Guild(ctx.Message.GuildID)

	msg := strings.Join(ctx.Args[1:], " ")

	embed := &discordgo.MessageEmbed{
		Description: msg,
		Color:       0x21a1a8,
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: guild.IconURL(),
			Name:    guild.Name + " Mod Team",
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Message ID: " + ctx.Message.ID,
		},
	}

	ch, err := ctx.Session.UserChannelCreate(m.User.ID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	send := &discordgo.MessageSend{
		Content: fmt.Sprintf("You have received a message from a moderator in **%v**:", guild.Name),
		Embed:   embed,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{},
		},
	}

	_, err = ctx.Session.ChannelMessageSendComplex(ch.ID, send)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	err = ctx.React(crouter.SuccessEmoji)
	return
}

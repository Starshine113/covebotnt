package crouter

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Send a message to the context channel
func (ctx *Ctx) Send(arg interface{}) (message *discordgo.Message, err error) {
	switch arg.(type) {
	case string:
		message, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, arg.(string))
	case *discordgo.MessageEmbed:
		message, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, arg.(*discordgo.MessageEmbed))
	case *discordgo.MessageSend:
		message, err = ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, arg.(*discordgo.MessageSend))
	default:
		err = errors.New("don't know what to do with that type")
	}
	return message, err
}

// Sendf sends a fmt.Sprintf-like input string
func (ctx *Ctx) Sendf(args ...interface{}) (msg *discordgo.Message, err error) {
	if len(args) < 2 {
		return nil, errors.New("not enough arguments")
	}
	format := args[0].(string)
	args = args[1:]
	msg, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf(format, args...))
	return
}

// Embed sends the input as an embed
func (ctx *Ctx) Embed(title, description string, color int) (msg *discordgo.Message, err error) {
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	msg, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
	return
}

// Edit a message
func (ctx *Ctx) Edit(message *discordgo.Message, arg interface{}) (msg *discordgo.Message, err error) {
	switch arg.(type) {
	case string:
		msg, err = ctx.Session.ChannelMessageEdit(message.ChannelID, message.ID, arg.(string))
	case *discordgo.MessageEmbed:
		msg, err = ctx.Session.ChannelMessageEditEmbed(message.ChannelID, message.ID, arg.(*discordgo.MessageEmbed))
	case *discordgo.MessageEdit:
		edit := arg.(*discordgo.MessageEdit)
		edit.ID = message.ID
		edit.Channel = message.ChannelID
		msg, err = ctx.Session.ChannelMessageEditComplex(edit)
	default:
		err = errors.New("don't know what to do with that type")
	}
	return msg, err
}

// React reacts to the triggering message
func (ctx *Ctx) React(emoji string) (err error) {
	return ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, emoji)
}

// TriggerTyping triggers typing in the channel the command was invoked in.
func (ctx *Ctx) TriggerTyping() (err error) {
	err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
	return
}

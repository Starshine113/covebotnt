package cbctx

import (
	"errors"

	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/bwmarrin/discordgo"
)

const (
	// SuccessEmoji is the emoji used to designate a successful action
	SuccessEmoji = "✅"
	// ErrorEmoji is the emoji used for errors
	ErrorEmoji = "❌"
	// WarnEmoji is the emoji used to warn that a command failed
	WarnEmoji = "⚠️"
)

// Ctx is the context for a command
type Ctx struct {
	Command          string
	Args             []string
	Session          *discordgo.Session
	Message          *discordgo.MessageCreate
	Channel          *discordgo.Channel
	Author           *discordgo.User
	BotUser          *discordgo.User
	Database         *cbdb.Db
	AdditionalParams map[string]interface{}
}

// Context creates a new Ctx
func Context(command string, args []string, s *discordgo.Session, m *discordgo.MessageCreate, db *cbdb.Db) (ctx *Ctx, err error) {
	ctx = &Ctx{Command: command, Args: args, Session: s, Message: m, Author: m.Author, Database: db}

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		return ctx, err
	}
	ctx.Channel = channel

	botUser, err := s.User("@me")
	if err != nil {
		return ctx, err
	}
	ctx.BotUser = botUser

	return ctx, nil
}

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

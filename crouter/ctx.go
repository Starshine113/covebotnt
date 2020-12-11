package crouter

import (
	"errors"
	"strings"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/Starshine113/covebotnt/bot"
	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/Starshine113/covebotnt/etc"
	"github.com/Starshine113/covebotnt/structs"
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
	GuildPrefix string
	Command     string
	Args        []string
	RawArgs     string

	Session  *discordgo.Session
	Bot      *bot.Bot
	BotUser  *discordgo.User
	Database *cbdb.Db
	BoltDb   *cbdb.BoltDb

	Message *discordgo.MessageCreate
	Channel *discordgo.Channel
	Author  *discordgo.User

	Handlers         *ttlcache.Cache
	AdditionalParams map[string]interface{}
	GuildSettings    *structs.GuildSettings
	Cmd              *Command

	Prefixes []string
}

// Errors when creating Context
var (
	ErrorNoBotUser = errors.New("bot user not found in state cache")
)

// Context creates a new Ctx
func Context(prefixes []string, messageContent string, m *discordgo.MessageCreate, b *bot.Bot) (ctx *Ctx, err error) {
	botUser := b.Session.State.User
	if botUser == nil {
		return nil, ErrorNoBotUser
	}

	var prefix string
	for _, p := range prefixes {
		if strings.HasPrefix(messageContent, p) {
			prefix = p
			break
		}
	}

	messageContent = etc.TrimPrefixesSpace(messageContent, prefixes...)
	message := strings.Split(messageContent, " ")
	command := message[0]
	args := []string{}
	if len(message) > 1 {
		args = message[1:]
	}

	ctx = &Ctx{GuildPrefix: prefix, Command: command, Args: args, Session: b.Session, Message: m, Author: m.Author, Database: b.Pool, BoltDb: b.Bolt, Handlers: b.Handlers, Bot: b, RawArgs: strings.Join(args, " "), Prefixes: prefixes}

	channel, err := b.Session.Channel(m.ChannelID)
	if err != nil {
		return ctx, err
	}
	ctx.Channel = channel
	ctx.BotUser = botUser

	return ctx, nil
}

package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// command handler
func messageCreateCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// if message was sent by a bot return; not only to ignore bots, but also to make sure PluralKit users don't trigger commands twice.
	if m.Author.Bot {
		return
	}

	// get prefix for the guild
	prefix := getPrefix(m.GuildID)

	// check if the message might be a command
	if strings.HasPrefix(strings.ToLower(m.Content), strings.ToLower(prefix)) {
		// split command(?) into a slice
		message := strings.Split(m.Content, " ")
		// remove prefix from command(?)
		message[0] = strings.TrimPrefix(strings.ToLower(message[0]), strings.ToLower(prefix))

		// get the command(?) and the args
		command := message[0]
		args := []string{}
		if len(message) > 1 {
			args = message[1:]
		}

		// run commandTree
		commandTree(command, args, s, m)
	}
}

package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/Starshine113/covebotnt/cbdb"
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
		var err error
		if len(message) > 1 {
			args = message[1:]
			args, err = combineQuotedItems(args)
			if err != nil {
				sugar.Errorf("Error parsing message %v: %v", m.ID, err)
			}
		}

		// run commandTree
		ctx, err := cbctx.Context(command, args, s, m, &cbdb.Db{db})
		if err != nil {
			sugar.Errorf("Error getting context: %v", err)
			return
		}
		commandTree(ctx)
		return
	}

	// if not, check if the message *started* with a bot mention
	botUser, err := s.User("@me")
	if err != nil {
		sugar.Errorf("Error fetching bot user: %v", err)
	}
	match, _ := regexp.MatchString(fmt.Sprintf("^<@!?%v>", botUser.ID), m.Content)
	if match {
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The current prefix is `%v`", prefix))
	}
}

func combineQuotedItems(in []string) (out []string, err error) {
	var matchedQuote bool
	var beginQuote int
	for i, item := range in {
		if matchedQuote {
			if strings.HasSuffix(item, "\"") {
				item = strings.Join(in[beginQuote:i+1], " ")
				item = strings.Trim(item, "\"")
				matchedQuote = false
				out = append(out, item)
			}
			if matchedQuote && i == len(in)-1 {
				err = errors.New("unexpected end of input")
			}
			continue
		}
		if strings.HasPrefix(item, "\"") {
			matchedQuote = true
			beginQuote = i
			continue
		}
		out = append(out, item)
	}
	return out, err
}

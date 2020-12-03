package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"codeberg.org/eviedelta/dwhook"
	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// command handler
func messageCreateCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	var err error

	// if message was sent by a bot return; not only to ignore bots, but also to make sure PluralKit users don't trigger commands twice.
	if m.Author.Bot {
		allowed := false
		for _, bot := range config.Bot.AllowedBots {
			if bot == m.Author.ID {
				allowed = true
			}
		}
		if allowed != true {
			return
		}
	}

	err = router.Respond(s, m)
	if err != nil {
		sugar.Errorf("Error sending autoresponse: %v", err)
	}

	// get prefix for the guild
	prefix := getPrefix(m.GuildID)

	ctx, err := crouter.Context(prefix, m.Content, s, m, pool, boltDb, handlerMap, b)
	if err != nil {
		sugar.Errorf("Error getting context: %v", err)
		return
	}
	// check if the message might be a command
	if ctx.MatchPrefix() {
		execute(ctx)
		return
	}

	// if not, check if the message contains a bot mention, and ends with "hello"
	content := strings.ToLower(m.Content)
	match, _ := regexp.MatchString(fmt.Sprintf("<@!?%v>.*hello$", s.State.User.ID), content)
	match2, _ := regexp.MatchString(fmt.Sprintf("%vhello$", regexp.QuoteMeta(prefix)), content)
	if match || match2 {
		ctx, err = crouter.Context("--", "--hello", s, m, pool, boltDb, handlerMap, b)
		if err != nil {
			sugar.Errorf("Error getting context: %v", err)
			return
		}
		execute(ctx)
		return
	}

	match, _ = regexp.MatchString(fmt.Sprintf("^<@!?%v>", s.State.User.ID), content)
	if match {
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The current prefix is `%v`", prefix))
		if err != nil {
			sugar.Errorf("Error sending message: %v", err)
			return
		}
		return
	}

	// no commands were triggered, so check if it's a DM, if not, return
	if m.GuildID != "" {
		return
	}
	if config.Bot.DMWebhook != "" && m.Author.ID != s.State.User.ID {
		for _, u := range config.Bot.BlockedUsers {
			if m.Author.ID == u {
				_, err = s.ChannelMessageSend(m.ChannelID, "You are blocked from DMing the bot. Please DM a bot admin if you think this is in error.")
				if err != nil {
					sugar.Errorf("Error sending message: %v", err)
				}
				return
			}
		}

		timestamp, err := discordgo.SnowflakeTimestamp(m.ID)
		if err != nil {
			sugar.Errorf("Error when getting timestamp for message %v: %v", m.ID, err)
		}

		var attachmentLink string
		if len(m.Attachments) > 0 {
			match, _ := regexp.MatchString("\\.(png|jpg|jpeg|gif|webp)$", m.Attachments[0].URL)
			if match {
				attachmentLink = m.Attachments[0].URL
			}
		}

		msg := dwhook.Message{
			Content:   fmt.Sprintf("> %v got DMed a message by **%v** (%v/%v):", s.State.User.Username, m.Author.String(), m.Author.Mention(), m.Author.ID),
			Username:  s.State.User.Username,
			AvatarURL: s.State.User.AvatarURL("256"),
			Embeds: []dwhook.Embed{{
				Author: dwhook.EmbedAuthor{
					Name:    m.Author.String(),
					IconURL: m.Author.AvatarURL("128"),
				},
				Color:       0x21a1a8,
				Description: m.Content,
				Footer: dwhook.EmbedFooter{
					Text: "Message ID: " + m.ID,
				},
				Timestamp: timestamp.UTC().Format(time.RFC3339),
				Image: dwhook.EmbedImage{
					URL: attachmentLink,
				},
			}},
		}

		dwhook.SendTo(config.Bot.DMWebhook, msg)
	}
}

func execute(ctx *crouter.Ctx) {
	guildSettings, err := pool.GetGuildSettings(ctx.Message.GuildID)
	if err != nil {
		sugar.Errorf("Error running command %v", err)
	}
	ctx.AdditionalParams = map[string]interface{}{"config": config, "botVer": botVersion, "gitVer": string(gitOut), "startTime": startTime}

	err = router.Execute(ctx, &guildSettings)
	if err != nil {
		sugar.Errorf("Error running command %v", err)
	}
}

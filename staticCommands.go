package main

import (
	"fmt"
	"time"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/bwmarrin/discordgo"
)

// Ping command: replies with latency and message edit time
func commandPing(ctx *cbctx.Ctx) (err error) {
	heartbeat := ctx.Session.HeartbeatLatency().String()

	// get current time
	cmdStart := time.Now()

	// send initial message
	message, err := ctx.Send("Pong!\nHeartbeat: " + heartbeat)
	if err != nil {
		return fmt.Errorf("Ping: %w", err)
	}

	// get time difference, edit message
	diff := time.Now().Sub(cmdStart).String()
	_, err = ctx.Edit(message, message.Content+"\nMessage latency: "+diff)
	return err
}

// Help shows the help pages
func commandHelp(ctx *cbctx.Ctx) (err error) {
	if len(ctx.Args) == 0 {
		embed := &discordgo.MessageEmbed{
			Title:       "CoveBotn't help",
			Description: "CoveBotn't is a general purpose bot, with ~~a gatekeeper, moderation commands, and~~ starboard functionality.",
			Color:       0x21a1a8,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Created by Starshine System (Starshine ☀✨#5000) | CoveBotn't v0.3",
			},
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Source code", Value: "CoveBotn't is licensed under the GNU AGPLv3. The source code can be found [here](https://github.com/Starshine113/covebotnt).", Inline: false},
				{Name: "Invite", Value: "Invite the bot with [this](" + ctx.AdditionalParams[0].(botConfig).Bot.Invite + ") link.", Inline: false},
				{Name: "Basic commands", Value: "`ping`: show the bot's latency\n`help`: show this help page", Inline: false},
			},
		}

		_, err := ctx.Send(embed)
		if err != nil {
			return fmt.Errorf("Help: %w", err)
		}
	}
	return nil
}

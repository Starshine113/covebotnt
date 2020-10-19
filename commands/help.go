package commands

import (
	"fmt"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/Starshine113/covebotnt/structs"
	"github.com/bwmarrin/discordgo"
)

// Help shows the help pages
func Help(ctx *cbctx.Ctx) (err error) {
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
				{Name: "Invite", Value: "Invite the bot with [this](" + ctx.AdditionalParams["config"].(structs.BotConfig).Bot.Invite + ") link.", Inline: false},
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

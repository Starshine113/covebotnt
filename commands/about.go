package commands

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/Starshine113/covebotnt/structs"
	"github.com/bwmarrin/discordgo"
)

// About shows some info about the bot
func About(ctx *cbctx.Ctx) (err error) {
	embed := &discordgo.MessageEmbed{
		Title:       "About CoveBotn't",
		Description: "CoveBotn't is a general purpose bot, with ~~a gatekeeper~~, moderation commands, and starboard functionality.",
		Color:       0x21a1a8,
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: ctx.BotUser.AvatarURL("256"),
			Text:    "Created with discordgo " + discordgo.VERSION,
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Source code",
				Value:  "CoveBotn't is licensed under the GNU AGPLv3. The source code can be found [here](https://github.com/Starshine113/covebotnt).",
				Inline: false,
			},
			{
				Name:   "Invite",
				Value:  "Invite the bot with [this](" + ctx.AdditionalParams["config"].(structs.BotConfig).Bot.Invite + ") link.",
				Inline: false,
			},
			{
				Name:   "Bot version",
				Value:  fmt.Sprintf("v%v-%v", ctx.AdditionalParams["botVer"].(string), ctx.AdditionalParams["gitVer"].(string)),
				Inline: true,
			},
			{
				Name:   "Go version",
				Value:  runtime.Version(),
				Inline: true,
			},
			{
				Name:   "discordgo version",
				Value:  discordgo.VERSION,
				Inline: true,
			},
			{
				Name:   "Author",
				Value:  "<@694563574386786314>",
				Inline: false,
			},
		},
	}

	_, err = ctx.Send(embed)
	if err != nil {
		return fmt.Errorf("Help: %w", err)
	}
	return nil
}

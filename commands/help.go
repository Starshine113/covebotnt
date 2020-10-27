package commands

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/Starshine113/covebotnt/structs"
	"github.com/bwmarrin/discordgo"
)

// Help shows the help pages
func Help(ctx *cbctx.Ctx) (err error) {
	botUser, err := ctx.Session.User("@me")
	if err != nil {
		return err
	}

	if len(ctx.Args) == 0 {
		embed := &discordgo.MessageEmbed{
			Title:       "CoveBotn't help",
			Description: "CoveBotn't is a general purpose bot, with ~~a gatekeeper~~, moderation commands, and starboard functionality.",
			Color:       0x21a1a8,
			Footer: &discordgo.MessageEmbedFooter{
				IconURL: botUser.AvatarURL("256"),
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
					Name:   "Basic commands",
					Value:  "`ping`: show the bot's latency\n`help`: show this help page\n`enlarge`: enlarge an emote\n`info`: get info about a user\n`serverinfo`: get info about a server\n`hello`: say hi to the bot",
					Inline: false,
				},
				{
					Name:   "Helper commands",
					Value:  "All these commands require a helper role.\n`setnote`: set a note for a user\n`delnote`: remove a note by ID\n`notes`: show notes for a user",
					Inline: false,
				},
				{
					Name:   "Mod commands",
					Value:  "All these commands require a mod role.\n`starboard`: manage the starboard\n`modroles`: manage mod roles\n`echo`: send a message as the bot\n`prefix`: change the bot prefix\n`export`: export this server's notes as JSON",
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

		_, err := ctx.Send(embed)
		if err != nil {
			return fmt.Errorf("Help: %w", err)
		}
	}
	return nil
}

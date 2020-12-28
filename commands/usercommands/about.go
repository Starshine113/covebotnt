package usercommands

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

func about(ctx *crouter.Ctx) (err error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	startTime := ctx.Bot.StartTime.UTC()
	botAuthor, err := ctx.Session.User("694563574386786314")
	if err != nil {
		return err
	}

	invite := ctx.Bot.Config.Bot.Invite
	if invite == "" {
		invite = fmt.Sprintf("[Invite link](%v)", ctx.Invite())
	}

	embed := &discordgo.MessageEmbed{
		Title: "About",
		Color: 0x21a1a8,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Made with discordgo %v", discordgo.VERSION),
			IconURL: "https://raw.githubusercontent.com/bwmarrin/discordgo/master/docs/img/discordgo.png",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: ctx.BotUser.AvatarURL("256"),
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Bot version",
				Value:  fmt.Sprintf("v%v-%v", ctx.Bot.Version, ctx.Bot.GitVer),
				Inline: true,
			},
			{
				Name:   "discordgo version",
				Value:  fmt.Sprintf("%v (%v)", discordgo.VERSION, runtime.Version()),
				Inline: true,
			},
			{
				Name:   "Author",
				Value:  botAuthor.Mention() + " / " + botAuthor.String(),
				Inline: false,
			},
			{
				Name:   "Uptime",
				Value:  fmt.Sprintf("%v (since %v)", PrettyDurationString(time.Since(startTime)), startTime.Format("Jan _2 2006, 15:04:05 MST")),
				Inline: false,
			},
			{
				Name:   "Memory usage",
				Value:  fmt.Sprintf("%v MB", m.TotalAlloc/1024/1024),
				Inline: true,
			},
			{
				Name:   "Invite",
				Value:  invite,
				Inline: true,
			},
			{
				Name:   "Source code",
				Value:  "[GitHub](https://github.com/Starshine113/covebotnt) / Licensed under the [GNU AGPLv3](https://www.gnu.org/licenses/agpl-3.0.html)",
				Inline: false,
			},
		},
	}

	_, err = ctx.Send(embed)
	return
}

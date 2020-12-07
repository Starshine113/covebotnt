package triggers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

func show(ctx *crouter.Ctx) (err error) {
	guild, err := ctx.Session.Guild(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	triggers, err := ctx.Database.Triggers(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	if len(triggers) == 0 {
		_, err = ctx.SendfNoAddXHandler("%v doesn't have any registered triggers.", guild.Name)
		return
	}

	if len(ctx.Args) > 0 {
		// try parsing an int
		if i, err := strconv.Atoi(ctx.Args[0]); err == nil {
			fmt.Println(i)
			for _, t := range triggers {
				if t.ID == i {
					_, err = ctx.Send(triggerEmbed(t))
					return err
				}
			}
		}
		// it's not an int, so...
		for _, t := range triggers {
			if strings.ToLower(strings.Join(ctx.Args, " ")) == strings.ToLower(t.Trigger) {
				_, err = ctx.Send(triggerEmbed(t))
				return err
			}
		}
	}

	triggerSlices := make([][]*cbdb.Trigger, 0)

	for i := 0; i < len(triggers); i += 10 {
		end := i + 10

		if end > len(triggers) {
			end = len(triggers)
		}

		triggerSlices = append(triggerSlices, triggers[i:end])
	}

	embeds := make([]*discordgo.MessageEmbed, 0)

	for i, s := range triggerSlices {
		x := make([]string, 0)
		for _, t := range s {
			if t == nil {
				continue
			}
			x = append(x, fmt.Sprintf("%v: `%v`", t.ID, t.Trigger))
		}
		embeds = append(embeds, &discordgo.MessageEmbed{
			Title: "Triggers for " + guild.Name,
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Page %v/%v", i+1, len(triggerSlices)),
			},
			Description: strings.Join(x, "\n"),
			Timestamp:   time.Now().Format(time.RFC3339),
			Color:       0x21a1a8,
		})
	}

	if len(embeds) == 1 {
		_, err = ctx.Send(embeds[0])
		return
	}

	_, err = ctx.PagedEmbed(embeds)
	return
}

func triggerEmbed(t *cbdb.Trigger) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "```" + t.Trigger + "```",
		Description: t.Response,
		Color:       0x21a1a8,
		Timestamp:   t.Modified.UTC().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("ID: %v", t.ID),
		},
	}
}

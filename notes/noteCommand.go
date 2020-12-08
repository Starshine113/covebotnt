package notes

import (
	"fmt"
	"time"

	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// CommandNotes list all notes for a given user
func CommandNotes(ctx *crouter.Ctx) (err error) {
	if len(ctx.Args) < 1 {
		ctx.CommandError(&crouter.ErrorMissingRequiredArgs{
			RequiredArgs: "user: id/mention",
			MissingArgs:  "user: id/mention",
		})
		return nil
	}

	err = ctx.TriggerTyping()
	if err != nil {
		return err
	}

	var user *discordgo.User

	member, err := ctx.ParseMember(ctx.Args[0])
	if err == nil {
		user = member.User
	} else {
		user, err = ctx.Session.User(ctx.Args[0])
		if err != nil {
			ctx.CommandError(err)
			return nil
		}
	}

	notes, err := ctx.Database.Notes(ctx.Message.GuildID, user.ID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	if len(notes) == 0 {
		_, err = ctx.SendfNoAddXHandler("**%v** has no log entries.", user.String())
		return err
	}

	notes = reverseNotes(notes)

	if len(notes) == 0 {
		_, err = ctx.Send(fmt.Sprintf("User **%v** has no notes.", user.String()))
		return
	}

	noteSlices := make([][]*cbdb.Note, 0)

	for i := 0; i < len(notes); i += 5 {
		end := i + 5

		if end > len(notes) {
			end = len(notes)
		}

		noteSlices = append(noteSlices, notes[i:end])
	}

	embeds := make([]*discordgo.MessageEmbed, 0)

	for i, s := range noteSlices {
		x := make([]*discordgo.MessageEmbedField, 0)
		for _, t := range s {
			if t == nil {
				continue
			}
			x = append(x, noteField(ctx, t))
		}
		embeds = append(embeds, &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    user.String(),
				IconURL: user.AvatarURL("128"),
			},
			Title: "Notes",
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Page %v/%v | User ID: %v", i+1, len(noteSlices), user.ID),
			},
			Fields:    x,
			Timestamp: time.Now().Format(time.RFC3339),
			Color:     0x21a1a8,
		})
	}

	if len(embeds) == 1 {
		_, err = ctx.Send(embeds[0])
		return
	}
	_, err = ctx.PagedEmbed(embeds)
	return
}

func reverseNotes(s []*cbdb.Note) []*cbdb.Note {
	a := make([]*cbdb.Note, len(s))
	copy(a, s)

	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}

	return a
}

package notes

import (
	"fmt"
	"strconv"

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
	var page int64 = 0

	if len(ctx.Args) > 1 {
		page, err = strconv.ParseInt(ctx.Args[1], 0, 0)
		if err != nil {
			ctx.CommandError(err)
			return nil
		}
		page = page - 1
		if page < 0 {
			page = 0
		}
	}

	msg, err := ctx.Send("Working, please wait...")
	if err != nil {
		return err
	}
	err = ctx.TriggerTyping()
	if err != nil {
		return err
	}

	member, err := ctx.ParseMember(ctx.Args[0])
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	notes, err := ctx.Database.Notes(ctx.Message.GuildID, member.User.ID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	notes = reverseNotes(notes)

	if len(notes) == 0 {
		_, err = ctx.Send(fmt.Sprintf("User **%v** has no notes.", member.User.String()))
		return
	}

	noteFields := make([]*discordgo.MessageEmbedField, 0)

	minRange := page * 10
	maxRange := ((page + 1) * 10)
	if int64(len(notes)) < maxRange {
		maxRange = int64(len(notes))
	}
	if int64(len(notes)) < minRange {
		minRange = int64(len(notes)) - 10
		if minRange < 0 {
			minRange = 0
		}
	}
	noteSlice := notes[minRange:maxRange]

	for _, note := range noteSlice {
		field, err := noteField(ctx.Session, note)
		if err != nil {
			return err
		}
		noteFields = append(noteFields, field)
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    member.User.String(),
			IconURL: member.User.AvatarURL("256"),
		},
		Title: "Notes",
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%v-%v out of %v shown | ID: %v", minRange+1, minRange+int64(len(noteSlice)), len(notes), member.User.ID),
		},
		Fields:    noteFields,
		Timestamp: string(ctx.Message.Timestamp),
	}

	_, err = ctx.Send(embed)
	if err != nil {
		return err
	}
	err = ctx.Session.ChannelMessageDelete(msg.ChannelID, msg.ID)
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

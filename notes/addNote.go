package notes

import (
	"fmt"
	"strings"
	"time"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/bwmarrin/discordgo"
)

// CommandSetNote sets a note
func CommandSetNote(ctx *cbctx.Ctx) (err error) {
	if len(ctx.Args) <= 1 {
		ctx.CommandError(&cbctx.ErrorMissingRequiredArgs{
			RequiredArgs: "user: id/mention, note: string",
			MissingArgs:  "user: id/mention, note: string",
		})
		return nil
	}

	member, err := ctx.ParseMember(ctx.Args[0])
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	note := strings.Join(ctx.Args[1:], " ")
	if len(note) > 200 {
		_, err = ctx.Send(fmt.Sprintf("%v This note is too long (maximum 200 characters). Input was %v characters.", cbctx.WarnEmoji, len(note)))
		return
	}

	err = ctx.Database.AddNote(&cbdb.Note{
		GuildID: ctx.Message.GuildID,
		UserID:  member.User.ID,
		ModID:   ctx.Author.ID,
		Note:    note,
	})
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	_, err = ctx.Send(fmt.Sprintf("%v ***Note taken.***\n**Note:** %v", cbctx.SuccessEmoji, note))
	return nil
}

func noteField(s *discordgo.Session, note *cbdb.Note) (field *discordgo.MessageEmbedField, err error) {
	mod, err := s.State.Member(note.GuildID, note.ModID)
	if err == discordgo.ErrStateNotFound {
		mod, err = s.GuildMember(note.GuildID, note.ModID)
	}
	var noteHeader string
	if err == nil {
		noteHeader = fmt.Sprintf("#%v (%v)", note.ID, mod.User.String())
	} else {
		noteHeader = fmt.Sprintf("#%v (%v)", note.ID, note.ModID)
	}

	return &discordgo.MessageEmbedField{
		Name:   noteHeader,
		Value:  fmt.Sprintf("%v\nCreated at %v", note.Note, note.Created.Format(time.RFC1123)),
		Inline: false,
	}, nil
}

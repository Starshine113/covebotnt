package notes

import (
	"fmt"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/bwmarrin/discordgo"
)

// CommandNotes list all notes for a given user
func CommandNotes(ctx *cbctx.Ctx) (err error) {
	err = ctx.TriggerTyping()
	if err != nil {
		return err
	}

	if len(ctx.Args) != 1 {
		ctx.CommandError(&cbctx.ErrorMissingRequiredArgs{
			RequiredArgs: "user: id/mention",
			MissingArgs:  "user: id/mention",
		})
		return nil
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

	if len(notes) == 0 {
		_, err = ctx.Send(fmt.Sprintf("User **%v#%v** has no notes.", member.User.Username, member.User.Discriminator))
		return
	}

	noteFields := make([]*discordgo.MessageEmbedField, 0)

	for _, note := range notes {
		field, err := noteField(ctx.Session, note)
		if err != nil {
			return err
		}
		noteFields = append(noteFields, field)
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("%v#%v", member.User.Username, member.User.Discriminator),
			IconURL: member.User.AvatarURL("256"),
		},
		Title: "Notes",
		Footer: &discordgo.MessageEmbedFooter{
			Text: "ID: " + member.User.ID,
		},
		Fields:    noteFields,
		Timestamp: string(ctx.Message.Timestamp),
	}

	_, err = ctx.Send(embed)
	return
}

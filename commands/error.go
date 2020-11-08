package commands

import (
	"fmt"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

// Error gets an error from the database
func Error(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckRequiredArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	id, err := uuid.Parse(ctx.Args[0])
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	cmdErr, err := ctx.BoltDb.GetError(id.String())
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	if cmdErr.Error == "" {
		_, err = ctx.Send(&discordgo.MessageEmbed{
			Title:       "Error not found",
			Description: fmt.Sprintf("There is no error with ID `%v`.", id.String()),
			Color:       0xbf1122,
		})
		return err
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Error `%v`", cmdErr.ErrorID),
		Description: fmt.Sprintf("```%v```", cmdErr.Error),
		Color:       0xbf1122,
	}

	_, err = ctx.Send(embed)
	return
}

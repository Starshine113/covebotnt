package notes

import (
	"fmt"
	"strconv"

	"github.com/Starshine113/covebotnt/cbctx"
)

// CommandDelNote deletes a note by ID
func CommandDelNote(ctx *cbctx.Ctx) (err error) {
	err = ctx.TriggerTyping()
	if err != nil {
		return err
	}

	if len(ctx.Args) != 1 {
		ctx.CommandError(&cbctx.ErrorMissingRequiredArgs{
			RequiredArgs: "id: int",
			MissingArgs:  "id: int",
		})
		return nil
	}

	id, err := strconv.ParseInt(ctx.Args[0], 0, 0)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	err = ctx.Database.DelNote(ctx.Message.GuildID, int(id))
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	_, err = ctx.Send(fmt.Sprintf("%v Removed note #%v.", cbctx.SuccessEmoji, id))
	return
}

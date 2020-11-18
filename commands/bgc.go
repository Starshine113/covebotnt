package commands

import (
	usercommands "github.com/Starshine113/covebotnt/commands/usercommands"
	"github.com/Starshine113/covebotnt/crouter"
	"github.com/Starshine113/covebotnt/notes"
)

// BGC is a combination command that runs UserInfo, Notes, and ModLogs for the specified user
func BGC(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	if err = usercommands.UserInfo(ctx); err != nil {
		return err
	}
	if err = notes.CommandNotes(ctx); err != nil {
		return err
	}
	if err = ModLogs(ctx); err != nil {
		return err
	}
	return
}

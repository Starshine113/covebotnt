package modcommands

import (
	"github.com/Starshine113/covebotnt/crouter"
)

// BGC is a combination command that runs UserInfo, Notes, and ModLogs for the specified user
func BGC(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	if err = ctx.Cmd.Router.GetCommand("userinfo").Command(ctx); err != nil {
		return err
	}
	if err = ctx.Cmd.Router.GetCommand("notes").Command(ctx); err != nil {
		return err
	}
	if err = ctx.Cmd.Router.GetCommand("modlogs").Command(ctx); err != nil {
		return err
	}
	return
}

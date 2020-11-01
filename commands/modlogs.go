package commands

import (
	"strings"

	"github.com/Starshine113/covebotnt/cbctx"
)

// ModLogs shows the moderation logs for the specified user
func ModLogs(ctx *cbctx.Ctx) (err error) {
	if len(ctx.Args) < 1 {
		ctx.CommandError(&cbctx.ErrorNotEnoughArgs{
			NumRequiredArgs: 1,
			SuppliedArgs:    0,
		})
		return nil
	}

	member, err := ctx.ParseMember(strings.Join(ctx.Args, " "))
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	_, err = ctx.Database.GetModLogs(ctx.Message.GuildID, member.User.ID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	return
}

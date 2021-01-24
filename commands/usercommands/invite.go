package usercommands

import (
	"github.com/starshine-sys/covebotnt/crouter"
)

// Invite sends an invite link for the bot
func Invite(ctx *crouter.Ctx) (err error) {
	_, err = ctx.Sendf("Use this link to invite me to your server: <%v>", ctx.Invite())
	return
}

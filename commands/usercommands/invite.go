package usercommands

import (
	"github.com/starshine-sys/covebotnt/crouter"
)

// Invite sends an invite link for the bot
func Invite(ctx *crouter.Ctx) (err error) {
	if ctx.Bot.Config.Bot.Public {
		_, err = ctx.Sendf("Use this link to invite me to your server: <%v>", ctx.Invite())
		return
	}

	// hardcoding user IDs is a good idea right
	u, err := ctx.ParseUser("694563574386786314")
	if err != nil {
		return err
	}

	_, err = ctx.Sendf(`Hi! This bot is currently invite-only.
If you'd like to invite it to your server, please ask us (%s), or DM the bot directly (as non-command DMs get sent to us directly). Note that while the bot *is* technically public, in the end, it's still a personal bot made for a single server, so we might not let you add it somewhere.

Alternatively, you can grab the source code at <https://github.com/starshine-sys/covebotnt> and run an instance yourself 🙂`, u)
	return err
}

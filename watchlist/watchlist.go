package watchlist

import (
	"strings"

	"github.com/starshine-sys/covebotnt/cbdb"
	"github.com/starshine-sys/covebotnt/crouter"
)

func watchlistAdd(ctx *crouter.Ctx) (err error) {
	if ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	u, err := ctx.ParseUser(strings.Join(ctx.Args, " "))
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	err = ctx.Database.AddToWatchlist(ctx.Message.GuildID, u.ID)
	if err != nil {
		if err == cbdb.ErrorAlreadyOnWatchlist {
			_, err = ctx.Send("That user is already on the watchlist.")
			return err
		}

		_, err = ctx.CommandError(err)
		return err
	}

	_, err = ctx.Sendf("Added %v to the watchlist.", u.Mention())
	return
}

func watchlistRemove(ctx *crouter.Ctx) (err error) {
	if ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	u, err := ctx.ParseUser(strings.Join(ctx.Args, " "))
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	err = ctx.Database.RemoveFromWatchlist(ctx.Message.GuildID, u.ID)
	if err != nil {
		if err == cbdb.ErrorNotOnWatchlist {
			_, err = ctx.Send("That user isn't watchlisted.")
			return err
		}

		_, err = ctx.CommandError(err)
		return err
	}

	_, err = ctx.Sendf("Removed %v from the watchlist.", u.Mention())
	return
}

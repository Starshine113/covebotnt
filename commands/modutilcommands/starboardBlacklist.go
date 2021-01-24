package modutilcommands

import (
	"strings"

	"github.com/starshine-sys/covebotnt/crouter"
)

// SbBlacklist ...
func SbBlacklist(ctx *crouter.Ctx) (err error) {
	sb, err := ctx.Database.GetStarboardBlacklist(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	x := make([]string, 0)
	for _, c := range sb {
		x = append(x, "<#"+c+">")
	}

	_, err = ctx.Embed("Starboard channel blacklist", strings.Join(x, "\n"), 0)
	return err
}

// SbBlacklistAdd ...
func SbBlacklistAdd(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	if strings.Join(ctx.Args, " ") == "" {
		return nil
	}

	c, err := ctx.ParseChannel(strings.Join(ctx.Args, " "))
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	err = ctx.Database.AddToStarboardBlacklist(ctx.Message.GuildID, c.ID)
	if err != nil {
		_, err = ctx.CommandError(err)
	} else {
		_, err = ctx.Sendf("%v Added %v to the starboard blacklist.", crouter.SuccessEmoji, c.Mention())
	}

	return
}

// SbBlacklistRemove ...
func SbBlacklistRemove(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	if strings.Join(ctx.Args, " ") == "" {
		return nil
	}

	c, err := ctx.ParseChannel(strings.Join(ctx.Args, " "))
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	err = ctx.Database.RemoveFromStarboardBlacklist(ctx.Message.GuildID, c.ID)
	if err != nil {
		_, err = ctx.CommandError(err)
	} else {
		_, err = ctx.Sendf("%v Removed %v from the starboard blacklist.", crouter.SuccessEmoji, c.Mention())
	}

	return
}

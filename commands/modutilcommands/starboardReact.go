package modutilcommands

import "github.com/Starshine113/covebotnt/crouter"

// StarboardReact ...
func StarboardReact(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckRequiredArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	err = ctx.Database.StarboardEmoji(ctx.Message.GuildID, ctx.Args[0])
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	g, err := ctx.Database.GetGuildSettings(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	_, err = ctx.Embedf("Starboard emoji changed", "Changed the starboard emoji to %v.", g.Starboard.Emoji)
	return
}

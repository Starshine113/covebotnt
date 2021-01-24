package modutilcommands

import "github.com/starshine-sys/covebotnt/crouter"

// StarboardSenderCanReactToggle ...
func StarboardSenderCanReactToggle(ctx *crouter.Ctx) (err error) {
	err = ctx.Database.ToggleSenderCanReact(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	g, err := ctx.Database.GetGuildSettings(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	if g.Starboard.SenderCanReact {
		_, err = ctx.Embedf("Changed settings", "The sender of a message can now react to it with the server's starboard emoji.")
	} else {
		_, err = ctx.Embedf("Changed settings", "The sender of a message can no longer react to it with the server's starboard emoji.")
	}
	return
}

package admincommands

import (
	"context"
	"errors"

	"github.com/Starshine113/covebotnt/crouter"
)

// Prefix ...
func Prefix(ctx *crouter.Ctx) (err error) {
	gs, err := ctx.Database.GetGuildSettings(ctx.Message.GuildID)
	if err != nil {
		return err
	}

	// if there are no arguments, show the current prefix
	if len(ctx.Args) == 0 {
		if gs.Prefix == "" {
			_, err = ctx.Send("The current prefix is `" + ctx.Bot.Config.Bot.Prefixes[0] + "` (default).")
			if err != nil {
				return err
			}
		} else {
			_, err = ctx.Send("The current prefix is `" + gs.Prefix + "`.")
			if err != nil {
				return err
			}
		}
		return nil
	}

	// if there's more than 1 argument, error
	if err = ctx.CheckRequiredArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	// otherwise, set prefix to first argument
	err = setGuildPrefix(ctx, ctx.Args[0], ctx.Message.GuildID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	_, err = ctx.Send("Changed prefix to `" + ctx.Args[0] + "`.")
	if err != nil {
		return err
	}

	return nil
}

func setGuildPrefix(ctx *crouter.Ctx, prefix, guildID string) error {
	ctx.Bot.Sugar.Infof("Changing prefix for %v", guildID)
	commandTag, err := ctx.Database.Pool.Exec(context.Background(), "update public.guild_settings set prefix = $1 where guild_id = $2", prefix, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	err = ctx.Database.RemoveFromGuildCache(guildID)
	if err != nil {
		return err
	}
	ctx.Bot.Sugar.Infof("Refreshed the settings for %v", guildID)
	return nil
}

package main

import (
	"context"
	"errors"

	"github.com/Starshine113/covebotnt/crouter"
)

func commandPrefix(ctx *crouter.Ctx) (err error) {
	gs, err := pool.GetGuildSettings(ctx.Message.GuildID)
	if err != nil {
		return err
	}

	// if there are no arguments, show the current prefix
	if len(ctx.Args) == 0 {
		if gs.Prefix == "" {
			_, err = ctx.Send("The current prefix is `" + config.Bot.Prefixes[0] + "` (default).")
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
	if len(ctx.Args) > 1 {
		ctx.CommandError(&crouter.ErrorTooManyArguments{1, len(ctx.Args)})
		return nil
	}

	// otherwise, set prefix to first argument
	err = setGuildPrefix(ctx.Args[0], ctx.Message.GuildID)
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

func setGuildPrefix(prefix, guildID string) error {
	sugar.Infof("Changing prefix for %v", guildID)
	commandTag, err := pool.Pool.Exec(context.Background(), "update public.guild_settings set prefix = $1 where guild_id = $2", prefix, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	err = pool.RemoveFromGuildCache(guildID)
	if err != nil {
		return err
	}
	sugar.Infof("Refreshed the settings for %v", guildID)
	return nil
}

func getPrefix(guildID string) string {
	gs, err := pool.GetGuildSettings(guildID)
	if err != nil {
		return config.Bot.Prefixes[0]
	}

	if gs.Prefix != "" {
		return gs.Prefix
	}
	return config.Bot.Prefixes[0]
}

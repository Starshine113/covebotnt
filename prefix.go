package main

import (
	"context"
	"errors"

	"github.com/Starshine113/covebotnt/cbctx"
)

func commandPrefix(ctx *cbctx.Ctx) (err error) {
	err = ctx.TriggerTyping()
	if err != nil {
		return err
	}

	// if there are no arguments, show the current prefix
	if len(ctx.Args) == 0 {
		if globalSettings[ctx.Message.GuildID].Prefix == "" {
			_, err = ctx.Send("The current prefix is `" + config.Bot.Prefixes[0] + "` (default).")
			if err != nil {
				return err
			}
		} else {
			_, err = ctx.Send("The current prefix is `" + globalSettings[ctx.Message.GuildID].Prefix + "`.")
			if err != nil {
				return err
			}
		}
		return nil
	}

	// if there's more than 1 argument, error
	if len(ctx.Args) > 1 {
		ctx.CommandError(&cbctx.ErrorTooManyArguments{1, len(ctx.Args)})
		return nil
	}

	// otherwise, set prefix to first argument
	err = setGuildPrefix(ctx.Args[0], ctx.Message.GuildID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	_, err = ctx.Send("Changed prefix to `" + globalSettings[ctx.Message.GuildID].Prefix + "`.")
	if err != nil {
		return err
	}

	return nil
}

func setGuildPrefix(prefix, guildID string) error {
	sugar.Infof("Changing prefix for %v", guildID)
	commandTag, err := db.Exec(context.Background(), "update public.guild_settings set prefix = $1 where guild_id = $2", prefix, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	err = getSettingsForGuild(guildID)
	if err != nil {
		return err
	}
	sugar.Infof("Refreshed the settings for %v", guildID)
	return nil
}

func getPrefix(guildID string) string {
	if globalSettings[guildID].Prefix != "" {
		return globalSettings[guildID].Prefix
	}
	return config.Bot.Prefixes[0]
}

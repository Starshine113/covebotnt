package main

import (
	"context"
	"errors"
)

func setGuildPrefix(prefix, guildID string) error {
	sugar.Infof("Changing prefix for %v", guildID)
	commandTag, err := db.Exec(context.Background(), "update guild_settings set prefix = $1 where guild_id = $2", prefix, guildID)
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

package cbdb

import (
	"context"
	"errors"
)

// AddToStarboardBlacklist adds the given channelID to the blacklist for guildID
func (db *Db) AddToStarboardBlacklist(guildID, channelID string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.guild_settings set sb_blacklist = array_append(sb_blacklist, $1) where guild_id = $2", channelID, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	err = db.RemoveFromGuildCache(guildID)
	return err
}

// RemoveFromStarboardBlacklist removes the given channelID from the blacklist for guildID
func (db *Db) RemoveFromStarboardBlacklist(guildID, channelID string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.guild_settings set sb_blacklist = array_remove(sb_blacklist, $1) where guild_id = $2", channelID, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	err = db.RemoveFromGuildCache(guildID)
	return err
}

// InStarboardBlacklist checks if the given channelID is in the blacklist for guildID
func (db *Db) InStarboardBlacklist(guildID, channelID string) (b bool) {
	db.Pool.QueryRow(context.Background(), "select $1 = any(guild.sb_blacklist) from (select * from public.guild_settings where guild_id = $2) as guild", channelID, guildID).Scan(&b)
	return b
}

// GetStarboardBlacklist returns the channel blacklist for guildID
func (db *Db) GetStarboardBlacklist(guildID string) (b []string, err error) {
	err = db.Pool.QueryRow(context.Background(), "select sb_blacklist from public.guild_settings where guild_id = $1", guildID).Scan(&b)
	return b, err
}

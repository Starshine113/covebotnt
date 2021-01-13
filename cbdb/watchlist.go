package cbdb

import (
	"context"
	"errors"
)

// Errors for setting the watchlist
var (
	ErrorAlreadyOnWatchlist = errors.New("user is already on watchlist")
	ErrorNotOnWatchlist     = errors.New("user is not on watchlist")
	ErrorNoRowsAffected     = errors.New("no rows affected")
)

// SetWatchlistChannel sets the watchlist channel for the given guild
func (db *Db) SetWatchlistChannel(guildID, channelID string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.guild_settings set watchlist_channel = $1 where guild_id = $2", channelID, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	return nil
}

// OnWatchlist returns true if a user is on the watchlist
func (db *Db) OnWatchlist(guildID, userID string) (b bool) {
	db.Pool.QueryRow(context.Background(), "select $1 = any(server.watchlist) from (select * from public.guild_settings where guild_id = $2) as server", userID, guildID).Scan(&b)
	return b
}

// AddToWatchlist adds the given userID to the watchlist for guildID
func (db *Db) AddToWatchlist(guildID, userID string) (err error) {
	if db.OnWatchlist(guildID, userID) {
		return ErrorAlreadyOnWatchlist
	}
	commandTag, err := db.Pool.Exec(context.Background(), "update public.guild_settings set watchlist = array_append(watchlist, $1) where guild_id = $2", userID, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return err
}

// RemoveFromWatchlist removes the given userID from the watchlist for guildID
func (db *Db) RemoveFromWatchlist(guildID, userID string) (err error) {
	if !db.OnWatchlist(guildID, userID) {
		return ErrorNotOnWatchlist
	}
	commandTag, err := db.Pool.Exec(context.Background(), "update public.guild_settings set watchlist = array_remove(watchlist, $1) where guild_id = $2", userID, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return err
}

// GetWatchlist returns the watchlist for guildID
func (db *Db) GetWatchlist(guildID string) (b []string, err error) {
	err = db.Pool.QueryRow(context.Background(), "select watchlist from public.guild_settings where guild_id = $1", guildID).Scan(&b)
	return b, err
}

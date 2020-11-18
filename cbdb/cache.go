package cbdb

import (
	"context"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/Starshine113/covebotnt/structs"
)

// GetGuildSettings gets the guild settings for a specific guild
func (db *Db) GetGuildSettings(g string) (s structs.GuildSettings, err error) {
	setting, err := db.GuildCache.Get(g)
	if err == ttlcache.ErrNotFound {
		s, err = db.GetDBGuildSettings(g)
		if err != nil {
			return s, err
		}
		err = db.GuildCache.Set(g, s)
		return s, err
	}
	return setting.(structs.GuildSettings), err
}

// RemoveFromGuildCache removes a cache entry
func (db *Db) RemoveFromGuildCache(g string) (err error) {
	err = db.GuildCache.Remove(g)
	if err == ttlcache.ErrNotFound {
		return nil
	}
	return
}

// GetDBGuildSettings gets the guild settings from the database
func (db *Db) GetDBGuildSettings(g string) (s structs.GuildSettings, err error) {
	var (
		prefix string

		starboardChannel                 string
		reactLimit                       int
		emoji                            string
		senderCanReact, reactToStarboard bool

		modRoles, helperRoles       []string
		modLog, muteRole, pauseRole string

		gatekeeperRoles, memberRoles      []string
		gatekeeperChannel, welcomeChannel string
		gatekeeperMessage, welcomeMessage string
	)

	row := db.Pool.QueryRow(context.Background(), "select * from public.guild_settings where guild_id=$1", g)

	err = row.Scan(&g, &prefix, &starboardChannel, &reactLimit, &emoji, &senderCanReact, &reactToStarboard,
		&modRoles, &helperRoles, &modLog, &muteRole, &pauseRole,
		&gatekeeperRoles, &memberRoles, &gatekeeperChannel,
		&gatekeeperMessage, &welcomeChannel, &welcomeMessage)
	if err != nil {
		return s, err
	}

	s = structs.GuildSettings{
		Prefix:     prefix,
		Starboard:  structs.StarboardSettings{starboardChannel, reactLimit, emoji, senderCanReact, reactToStarboard},
		Moderation: structs.ModSettings{modRoles, helperRoles, modLog, muteRole, pauseRole},
		Gatekeeper: structs.GatekeeperSettings{gatekeeperRoles, memberRoles, gatekeeperChannel,
			gatekeeperMessage, welcomeChannel, welcomeMessage}}

	return s, nil
}

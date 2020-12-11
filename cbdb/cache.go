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
		prefixes []string

		starboardChannel                 string
		reactLimit                       int
		emoji                            string
		senderCanReact, reactToStarboard bool

		modRoles, helperRoles       []string
		modLog, muteRole, pauseRole string

		gatekeeperRoles, memberRoles      []string
		gatekeeperChannel, welcomeChannel string
		gatekeeperMessage, welcomeMessage string

		yagLog     string
		yagEnabled bool
	)

	row := db.Pool.QueryRow(context.Background(), `select
	g.prefixes, g.starboard_channel, g.react_limit,
	g.emoji, g.sender_can_react, g.react_to_starboard,
	g.mod_roles, g.helper_roles, g.mod_log, g.mute_role,
	g.pause_role, g.gatekeeper_roles, g.member_roles,
	g.gatekeeper_channel, g.gatekeeper_message,
	g.welcome_channel, g.welcome_message,
	y.log_channel, y.enabled
	from public.guild_settings as g, public.yag_import as y
	where g.guild_id=$1 and y.guild_id = $1`, g)

	err = row.Scan(&prefixes, &starboardChannel, &reactLimit, &emoji, &senderCanReact, &reactToStarboard,
		&modRoles, &helperRoles, &modLog, &muteRole, &pauseRole,
		&gatekeeperRoles, &memberRoles, &gatekeeperChannel,
		&gatekeeperMessage, &welcomeChannel, &welcomeMessage, &yagLog, &yagEnabled)
	if err != nil {
		return s, err
	}

	s = structs.GuildSettings{
		Prefixes: prefixes,
		Starboard: structs.StarboardSettings{
			StarboardID:      starboardChannel,
			ReactLimit:       reactLimit,
			Emoji:            emoji,
			SenderCanReact:   senderCanReact,
			ReactToStarboard: reactToStarboard,
		},
		Moderation: structs.ModSettings{
			ModRoles:    modRoles,
			HelperRoles: helperRoles,
			ModLog:      modLog,
			MuteRole:    muteRole,
			PauseRole:   pauseRole,
		},
		Gatekeeper: structs.GatekeeperSettings{
			GatekeeperRoles:   gatekeeperRoles,
			MemberRoles:       memberRoles,
			GatekeeperChannel: gatekeeperChannel,
			GatekeeperMessage: gatekeeperMessage,
			WelcomeChannel:    welcomeChannel,
			WelcomeMessage:    welcomeMessage,
		},
		YAG: structs.YAGImport{
			Channel: yagLog,
			Enabled: yagEnabled,
		},
	}

	return s, nil
}

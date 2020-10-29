package main

import (
	"context"

	"github.com/Starshine113/covebotnt/structs"
)

type guildSettings struct {
	Prefix     string
	Starboard  starboardSettings
	Moderation modSettings
	Gatekeeper gatekeeperSettings
}
type starboardSettings struct {
	StarboardID      string
	ReactLimit       int
	Emoji            string
	SenderCanReact   bool
	ReactToStarboard bool
}
type modSettings struct {
	ModRoles    []string
	HelperRoles []string
	ModLog      string
	MuteRole    string
	PauseRole   string
}
type gatekeeperSettings struct {
	GatekeeperRoles   []string
	MemberRoles       []string
	GatekeeperChannel string
	GatekeeperMessage string
	WelcomeChannel    string
	WelcomeMessage    string
}

func initSettingsForGuild(guildID string) (err error) {
	_, err = db.Exec(context.Background(), "insert into public.guild_settings (guild_id) values ($1)", guildID)
	if err != nil {
		return err
	}
	err = getSettingsForGuild(guildID)
	return err
}

func getSettingsForGuild(guildID string) (err error) {
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

	row := db.QueryRow(context.Background(), "select * from public.guild_settings where guild_id=$1", guildID)

	row.Scan(&guildID, &prefix, &starboardChannel, &reactLimit, &emoji, &senderCanReact, &reactToStarboard,
		&modRoles, &helperRoles, &modLog, &muteRole, &pauseRole,
		&gatekeeperRoles, &memberRoles, &gatekeeperChannel,
		&gatekeeperMessage, &welcomeChannel, &welcomeMessage)

	globalSettings[guildID] = structs.GuildSettings{
		Prefix:     prefix,
		Starboard:  structs.StarboardSettings{starboardChannel, reactLimit, emoji, senderCanReact, reactToStarboard},
		Moderation: structs.ModSettings{modRoles, helperRoles, modLog, muteRole, pauseRole},
		Gatekeeper: structs.GatekeeperSettings{gatekeeperRoles, memberRoles, gatekeeperChannel,
			gatekeeperMessage, welcomeChannel, welcomeMessage}}

	return nil
}

func getSettingsAll() (map[string]structs.GuildSettings, error) {
	settings := make(map[string]structs.GuildSettings)
	// get starboard settings
	rows, err := db.Query(context.Background(), "select * from public.guild_settings")
	if err != nil {
		return settings, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			guildID, prefix string

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
		rows.Scan(
			&guildID, &prefix, &starboardChannel, &reactLimit, &emoji, &senderCanReact, &reactToStarboard,
			&modRoles, &helperRoles, &modLog, &muteRole, &pauseRole,
			&gatekeeperRoles, &memberRoles, &gatekeeperChannel,
			&gatekeeperMessage, &welcomeChannel, &welcomeMessage,
		)
		settings[guildID] = structs.GuildSettings{
			Prefix:     prefix,
			Starboard:  structs.StarboardSettings{starboardChannel, reactLimit, emoji, senderCanReact, reactToStarboard},
			Moderation: structs.ModSettings{modRoles, helperRoles, modLog, muteRole, pauseRole},
			Gatekeeper: structs.GatekeeperSettings{gatekeeperRoles, memberRoles, gatekeeperChannel,
				gatekeeperMessage, welcomeChannel, welcomeMessage}}
	}

	return settings, err
}

package main

import (
	"context"
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

func updateSettingsForGuild(guildID string) error {
	return nil
}

func getSettingsAll() (settings map[string]guildSettings, err error) {
	// get starboard settings
	rows, err := db.Query(context.Background(), "select * from guild_settings")
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
		settings[guildID] = guildSettings{
			Prefix:     prefix,
			Starboard:  starboardSettings{starboardChannel, reactLimit, emoji, senderCanReact, reactToStarboard},
			Moderation: modSettings{modRoles, helperRoles, modLog, muteRole, pauseRole},
			Gatekeeper: gatekeeperSettings{gatekeeperRoles, memberRoles, gatekeeperChannel,
				gatekeeperMessage, welcomeChannel, welcomeMessage}}
	}

	return settings, err
}

func setStarboardChannel(guildID int, channelID int) error {
	return nil
}

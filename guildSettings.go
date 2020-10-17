package main

import (
	"context"
)

type guildSettings struct {
	Starboard  starboardSettings
	Moderation modSettings
	Gatekeeper gatekeeperSettings
}
type starboardSettings struct {
	StarboardID      uint64
	ReactLimit       int
	Emoji            string
	SenderCanReact   bool
	ReactToStarboard bool
}
type modSettings struct {
	ModRoles    []uint64
	HelperRoles []uint64
	ModLog      uint64
	MuteRole    uint64
	PauseRole   uint64
}
type gatekeeperSettings struct {
	GatekeeperRoles   []uint64
	MemberRoles       []uint64
	GatekeeperChannel uint64
	GatekeeperMessage string
	WelcomeChannel    uint64
	WelcomeMessage    string
}

func updateSettingsForGuild(guildID uint64) error {
	return nil
}

func getSettingsAll() (map[uint64]guildSettings, error) {
	settings := make(map[uint64]guildSettings)

	// get starboard settings
	rows, err := db.Query(context.Background(), "select * from guild_settings")
	if err != nil {
		return settings, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			guildID, starboardChannel        uint64
			reactLimit                       int
			emoji                            string
			senderCanReact, reactToStarboard bool

			modRoles, helperRoles       []uint64
			modLog, muteRole, pauseRole uint64

			gatekeeperRoles, memberRoles      []uint64
			gatekeeperChannel, welcomeChannel uint64
			gatekeeperMessage, welcomeMessage string
		)
		rows.Scan(
			&guildID, &starboardChannel, &reactLimit, &emoji, &senderCanReact, &reactToStarboard,
			&modRoles, &helperRoles, &modLog, &muteRole, &pauseRole,
			&gatekeeperRoles, &memberRoles, &gatekeeperChannel,
			&gatekeeperMessage, &welcomeChannel, &welcomeMessage,
		)
		settings[guildID] = guildSettings{
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

package cbdb

import (
	"context"
	"errors"
)

// SetGatekeeperChannel sets the gatekeeper channel for the given guild
func (db *Db) SetGatekeeperChannel(guildID, channelID string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.guild_settings set gatekeeper_channel = $1 where guild_id = $2", channelID, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	return nil
}

// SetWelcomeChannel sets the welcome channel for the given guild
func (db *Db) SetWelcomeChannel(guildID, channelID string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.guild_settings set welcome_channel = $1 where guild_id = $2", channelID, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	return nil
}

// SetGatekeeperMsg sets the gatekeeper message for the given guild
func (db *Db) SetGatekeeperMsg(guildID, msg string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.guild_settings set gatekeeper_message = $1 where guild_id = $2", msg, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	return nil
}

// SetWelcomeMsg sets the welcome message for the given guild
func (db *Db) SetWelcomeMsg(guildID, msg string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.guild_settings set welcome_message = $1 where guild_id = $2", msg, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	return nil
}

// SetGatekeeperRoles sets the gatekeeper roles for the given guild
func (db *Db) SetGatekeeperRoles(guildID string, roles []string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.guild_settings set gatekeeper_roles = $1 where guild_id = $2", roles, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	return nil
}

// SetMemberRoles sets the gatekeeper roles for the given guild
func (db *Db) SetMemberRoles(guildID string, roles []string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.guild_settings set member_roles = $1 where guild_id = $2", roles, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	return nil
}

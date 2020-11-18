package cbdb

import "context"

// InitSettingsForGuild initialises the settings for a guild iif it doesn't have any yet
func (db *Db) InitSettingsForGuild(guildID string) (err error) {
	var exists bool
	err = db.Pool.QueryRow(context.Background(), "select exists (select from public.guild_settings where guild_id = $1)", guildID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	_, err = db.Pool.Exec(context.Background(), "insert into public.guild_settings (guild_id) values ($1)", guildID)
	return err
}

package cbdb

import (
	"context"
	"errors"
)

// InsertStarboardEntry inserts an entry into the database
func (db *Db) InsertStarboardEntry(messageID, channelID, guildID, starboardMessageID string) error {
	_, err := db.Pool.Exec(context.Background(), "insert into public.starboard_messages (message_id, channel_id, server_id, starboard_message_id) values ($1, $2, $3, $4)", messageID, channelID, guildID, starboardMessageID)
	return err
}

// DeleteStarboardEntry deletes an entry from the database
func (db *Db) DeleteStarboardEntry(messageID string) error {
	commandTag, err := db.Pool.Exec(context.Background(), "delete from public.starboard_messages where message_id = $1 or starboard_message_id = $2", messageID, messageID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}

	return nil
}

// SetStarboardChannel sets the starboard channel for a guild
func (db *Db) SetStarboardChannel(channelID, guildID string) error {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.guild_settings set starboard_channel = $1 where guild_id = $2", channelID, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	err = db.RemoveFromGuildCache(guildID)
	if err != nil {
		return err
	}
	return nil
}

// SetStarboardLimit sets the starboard limit for a guild
func (db *Db) SetStarboardLimit(limit int, guildID string) error {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.guild_settings set react_limit = $1 where guild_id = $2", limit, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	err = db.RemoveFromGuildCache(guildID)
	if err != nil {
		return err
	}
	return nil
}

// GetStarboardEntry gets the starboard entry for the given ID
func (db *Db) GetStarboardEntry(m string) (s string) {
	db.Pool.QueryRow(context.Background(), "select starboard_message_id from public.starboard_messages where message_id = $1", m).Scan(&s)
	return s
}

// GetOrigStarboardMessage ...
func (db *Db) GetOrigStarboardMessage(m string) (s string) {
	db.Pool.QueryRow(context.Background(), "select starboard_message_id from public.starboard_messages where starboard_message_id = $1", m).Scan(&s)
	return s
}

// ToggleSenderCanReact ...
func (db *Db) ToggleSenderCanReact(g string) (err error) {
	s, err := db.GetGuildSettings(g)
	if err != nil {
		return err
	}

	var b bool
	if s.Starboard.SenderCanReact {
		b = false
	} else {
		b = true
	}

	commandTag, err := db.Pool.Exec(context.Background(), "update public.guild_settings set sender_can_react = $1 where guild_id = $2", b, g)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	err = db.RemoveFromGuildCache(g)
	if err != nil {
		return err
	}
	return nil
}

// StarboardEmoji ...
func (db *Db) StarboardEmoji(g, e string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.guild_settings set emoji = $1 where guild_id = $2", e, g)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	err = db.RemoveFromGuildCache(g)
	if err != nil {
		return err
	}
	return nil
}

package cbdb

import (
	"context"
	"time"
)

// ModLogEntry is an entry in the mod log
type ModLogEntry struct {
	ID      int       `json:"id"`
	GuildID string    `json:"guild_id"`
	UserID  string    `json:"user_id"`
	ModID   string    `json:"mod_id"`
	Type    string    `json:"type"`
	Reason  string    `json:"reason"`
	Time    time.Time `json:"timestamp"`
}

// SetModLogChannel sets the moderation log channel for the guild
func (db *Db) SetModLogChannel(guildID, channelID string) (err error) {
	_, err = db.Pool.Exec(context.Background(), "update public.guild_settings set mod_log = $1 where guild_id = $2", channelID, guildID)
	return err
}

// AddToModLog adds the specified partial ModLogEntry to the moderation log, and returns the full object
func (db *Db) AddToModLog(entry *ModLogEntry) (out *ModLogEntry, err error) {
	var id int
	var timestamp time.Time

	err = db.Pool.QueryRow(context.Background(), "insert into public.mod_log (guild_id, user_id, mod_id, type, reason, created) values ($1, $2, $3, $4, $5, $6) returning id, created", entry.GuildID, entry.UserID, entry.ModID, entry.Type, entry.Reason, entry.Time).Scan(&id, &timestamp)
	if err != nil {
		return
	}
	out = entry
	out.ID = id
	out.Time = timestamp
	return
}

// GetModLogs gets all the mod logs for a user
func (db *Db) GetModLogs(guildID, userID string) (out []*ModLogEntry, err error) {
	rows, err := db.Pool.Query(context.Background(), "select id, guild_id, user_id, mod_id, type, reason, created from public.mod_log where user_id = $1 and guild_id = $2", userID, guildID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id                                         int
			guildID, userID, modID, actionType, reason string
			created                                    time.Time
		)

		rows.Scan(&id, &guildID, &userID, &modID, &actionType, &reason, &created)
		out = append(out, &ModLogEntry{
			ID:      id,
			GuildID: guildID,
			UserID:  userID,
			ModID:   modID,
			Type:    actionType,
			Reason:  reason,
			Time:    created,
		})
	}
	return
}

// GetAllLogs gets *all* the mod logs for a guild
func (db *Db) GetAllLogs(guildID string) (out []*ModLogEntry, err error) {
	rows, err := db.Pool.Query(context.Background(), "select id, guild_id, user_id, mod_id, type, reason, created from public.mod_log where guild_id = $1", guildID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id                                         int
			guildID, userID, modID, actionType, reason string
			created                                    time.Time
		)

		rows.Scan(&id, &guildID, &userID, &modID, &actionType, &reason, &created)
		out = append(out, &ModLogEntry{
			ID:      id,
			GuildID: guildID,
			UserID:  userID,
			ModID:   modID,
			Type:    actionType,
			Reason:  reason,
			Time:    created,
		})
	}
	return
}

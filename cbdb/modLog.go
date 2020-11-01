package cbdb

import (
	"context"
	"time"
)

// ModLogEntry is an entry in the mod log
type ModLogEntry struct {
	ID      int
	GuildID string
	UserID  string
	ModID   string
	Type    string
	Reason  string
	Time    time.Time
}

// SetModLogChannel sets the moderation log channel for the guild
func (db *Db) SetModLogChannel(guildID, channelID string) (err error) {
	_, err = db.Pool.Exec(context.Background(), "update public.guild_settings set mod_log = $1 where guild_id = $2", channelID, guildID)
	return err
}

// AddToModLog adds the specified partial ModLogEntry to the moderation log, and returns the full object
func (db *Db) AddToModLog(entry *ModLogEntry) (out *ModLogEntry, err error) {
	var id int

	err = db.Pool.QueryRow(context.Background(), "insert into public.mod_log (guild_id, user_id, mod_id, type, reason) values ($1, $2, $3, $4, $5) returning id", entry.GuildID, entry.UserID, entry.ModID, entry.Type, entry.Reason).Scan(&id)
	if err != nil {
		return
	}

	var timestamp time.Time

	err = db.Pool.QueryRow(context.Background(), "select created from public.mod_log where id=$1 ", id).Scan(&timestamp)
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
	rows, err := db.Pool.Query(context.Background(), "select * from public.mod_log where user_id = $1 and guild_id = $2", userID, guildID)
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

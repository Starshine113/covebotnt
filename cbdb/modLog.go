package cbdb

import (
	"context"
	"errors"
	"time"

	"github.com/Starshine113/snowflake"
)

// ModLogEntry is an entry in the mod log
type ModLogEntry struct {
	ID        int                 `json:"id"`
	GuildID   string              `json:"guild_id"`
	UserID    string              `json:"user_id"`
	ModID     string              `json:"mod_id"`
	Type      string              `json:"type"`
	Reason    string              `json:"reason"`
	Time      time.Time           `json:"timestamp"`
	Snowflake snowflake.Snowflake `json:"snowflake"`
}

// SetModLogChannel sets the moderation log channel for the guild
func (db *Db) SetModLogChannel(guildID, channelID string) (err error) {
	_, err = db.Pool.Exec(context.Background(), "update public.guild_settings set mod_log = $1 where guild_id = $2", channelID, guildID)
	return err
}

// General errors related to mod log operations
var (
	ErrSnowflakeAlreadyExists = errors.New("modlog: snowflake already exists in mod log")
)

// AddToModLog adds the specified partial ModLogEntry to the moderation log, and returns the full object
func (db *Db) AddToModLog(entry *ModLogEntry) (out *ModLogEntry, err error) {
	var id int
	var timestamp time.Time

	var snowflakeExists bool

	err = db.Pool.QueryRow(context.Background(), "select exists (select from public.mod_log where snowflake = $1)", entry.Snowflake).Scan(&snowflakeExists)
	if err != nil {
		return
	}

	if snowflakeExists {
		return entry, ErrSnowflakeAlreadyExists
	}

	err = db.Pool.QueryRow(context.Background(), "insert into public.mod_log (guild_id, user_id, mod_id, type, reason, created, snowflake) values ($1, $2, $3, $4, $5, $6, $7) returning id, created;", entry.GuildID, entry.UserID, entry.ModID, entry.Type, entry.Reason, entry.Time, entry.Snowflake).Scan(&id, &timestamp)
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
	rows, err := db.Pool.Query(context.Background(), "select id, guild_id, user_id, mod_id, type, reason, created, snowflake from public.mod_log where user_id = $1 and guild_id = $2 order by created desc", userID, guildID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id        int
			snowflake snowflake.Snowflake
			created   time.Time

			guildID, userID, modID, actionType, reason string
		)

		rows.Scan(&id, &guildID, &userID, &modID, &actionType, &reason, &created, &snowflake)
		out = append(out, &ModLogEntry{
			ID:        id,
			GuildID:   guildID,
			UserID:    userID,
			ModID:     modID,
			Type:      actionType,
			Reason:    reason,
			Time:      created,
			Snowflake: snowflake,
		})
	}
	return
}

// GetAllLogs gets *all* the mod logs for a guild
func (db *Db) GetAllLogs(guildID string) (out []*ModLogEntry, err error) {
	rows, err := db.Pool.Query(context.Background(), "select id, guild_id, user_id, mod_id, type, reason, created, snowflake from public.mod_log where guild_id = $1", guildID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id        int
			snowflake snowflake.Snowflake
			created   time.Time

			guildID, userID, modID, actionType, reason string
		)

		rows.Scan(&id, &guildID, &userID, &modID, &actionType, &reason, &created, &snowflake)
		out = append(out, &ModLogEntry{
			ID:        id,
			GuildID:   guildID,
			UserID:    userID,
			ModID:     modID,
			Type:      actionType,
			Reason:    reason,
			Time:      created,
			Snowflake: snowflake,
		})
	}
	return
}

package cbdb

import (
	"context"
	"errors"
	"time"

	"github.com/Starshine113/snowflake"
)

// Trigger ...
type Trigger struct {
	ID        int
	GuildID   string
	Creator   string
	Modified  time.Time
	Trigger   string
	Response  string
	Snowflake snowflake.Snowflake
}

// AddTrigger ...
func (db *Db) AddTrigger(t *Trigger) (*Trigger, error) {
	if t == nil {
		return nil, errors.New("trigger was nil")
	}
	if t.GuildID == "" || t.Creator == "" || t.Trigger == "" || t.Response == "" {
		return nil, errors.New("one or more required fields was nil")
	}
	if len(t.Trigger) > 99 {
		t.Trigger = t.Trigger[:99]
	}
	if len(t.Response) > 1999 {
		t.Response = t.Response[:1999]
	}
	var timestamp time.Time
	var id int

	err := db.Pool.QueryRow(context.Background(), "insert into public.triggers (guild_id, created_by, trigger, response, snowflake) values ($1, $2, $3, $4, $5) returning id, modified", t.GuildID, t.Creator, t.Trigger, t.Response, t.Snowflake).Scan(&id, &timestamp)
	t.ID = id
	t.Modified = timestamp
	return t, err
}

// EditTrigger ...
func (db *Db) EditTrigger(guildID string, triggerID int, t *Trigger) (*Trigger, error) {
	if t == nil {
		return nil, errors.New("trigger was nil")
	}
	if t.Trigger == "" || t.Response == "" {
		return nil, errors.New("one or more required fields was nil")
	}
	if len(t.Trigger) > 99 {
		t.Trigger = t.Trigger[:99]
	}
	if len(t.Response) > 1999 {
		t.Response = t.Response[:1999]
	}

	var modified time.Time

	err := db.Pool.QueryRow(context.Background(), "update public.triggers set trigger = $1, response = $2, modified = $3 where guild_id = $4 and id = $5 returning modified", t.Trigger, t.Response, time.Now().UTC(), guildID, triggerID).Scan(&modified)
	if err != nil {
		return t, err
	}
	t.Modified = modified
	return t, nil
}

// Triggers gets all triggers for a guild
func (db *Db) Triggers(id string) (out []*Trigger, err error) {
	rows, err := db.Pool.Query(context.Background(), "select id, guild_id, created_by, modified, trigger, response, snowflake from public.triggers where guild_id = $1 order by id asc", id)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id        int
			snowflake snowflake.Snowflake
			modified  time.Time

			guildID, creator, trigger, response string
		)

		rows.Scan(&id, &guildID, &creator, &modified, &trigger, &response, &snowflake)
		out = append(out, &Trigger{
			ID:        id,
			GuildID:   guildID,
			Creator:   creator,
			Modified:  modified,
			Trigger:   trigger,
			Response:  response,
			Snowflake: snowflake,
		})
	}

	return
}

// RemoveTrigger ...
func (db *Db) RemoveTrigger(guildID string, triggerID int) (err error) {
	_, err = db.Pool.Exec(context.Background(), "delete from public.triggers where id = $1 and guild_id = $2", triggerID, guildID)
	return
}

package cbdb

import (
	"context"
	"time"
)

// Note holds the data for a note
type Note struct {
	ID      int
	GuildID string
	UserID  string
	ModID   string
	Note    string
	Created time.Time
}

// AddNote adds a note to the database
func (db *Db) AddNote(note *Note) (err error) {
	_, err = db.Pool.Exec(context.Background(), "insert into notes (guild_id, user_id, mod_id, note) values ($1, $2, $3, $4)", note.GuildID, note.UserID, note.ModID, note.Note)
	return err
}

// DelNote deletes a note from the database
func (db *Db) DelNote(guildID string, id int) (err error) {
	_, err = db.Pool.Exec(context.Background(), "delete from notes where id = $1 and guild_id = $2", id, guildID)
	return
}

// Notes gets all the notes for a user in a guild
func (db *Db) Notes(guildID, userID string) (notes []*Note, err error) {
	rows, err := db.Pool.Query(context.Background(), "select * from notes where user_id = $1 and guild_id = $2", userID, guildID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id                           int
			guildID, userID, modID, note string
			created                      time.Time
		)

		rows.Scan(&id, &guildID, &userID, &modID, &note, &created)
		notes = append(notes, &Note{
			ID:      id,
			GuildID: guildID,
			UserID:  userID,
			ModID:   modID,
			Note:    note,
			Created: created,
		})
	}
	return
}

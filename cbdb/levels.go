package cbdb

import (
	"fmt"
	"sort"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

// GetXPForUser gets the user's XP
func (db *BoltDb) GetXPForUser(userID, guildID string) (xp int64, err error) {
	err = db.Bolt.Update(func(tx *bolt.Tx) (err error) {
		level := tx.Bucket([]byte("levels"))
		guild := level.Bucket([]byte(guildID))

		xpB := guild.Get([]byte(userID))
		if xpB == nil {
			err = guild.Put([]byte(userID), []byte("0"))
			return
		}

		xp, err = strconv.ParseInt(string(xpB), 10, 0)
		return
	})

	return
}

// AddXPForUser increments the user's XP by one
func (db *BoltDb) AddXPForUser(userID, guildID string) (newXp int64, err error) {
	err = db.Bolt.Update(func(tx *bolt.Tx) (err error) {
		level := tx.Bucket([]byte("levels"))
		guild := level.Bucket([]byte(guildID))

		xpB := guild.Get([]byte(userID))
		if xpB == nil {
			err = guild.Put([]byte(userID), []byte("0"))
			if err != nil {
				return
			}
			xpB = []byte("0")
		}

		xp, err := strconv.ParseInt(string(xpB), 10, 0)
		newXp = xp + 1
		err = guild.Put([]byte(userID), []byte(fmt.Sprintf("%v", newXp)))

		return
	})

	return
}

// Xp holds XP for a given user
type Xp struct {
	UserID string
	Xp     int64
}

// AllEntriesForGuild gets all the entries for a guild
func (db *BoltDb) AllEntriesForGuild(guildID string) (xp XpList, err error) {
	err = db.Bolt.View(func(tx *bolt.Tx) (err error) {
		level := tx.Bucket([]byte("levels"))
		guild := level.Bucket([]byte(guildID))

		c := guild.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			userXp, err := strconv.ParseInt(string(v), 10, 0)
			if err != nil {
				return err
			}
			xp = append(xp, &Xp{UserID: string(k), Xp: userXp})
		}

		return nil
	})

	sort.Sort(xp)
	return
}

// XpList is a slice of all XP in a guild
type XpList []*Xp

func (a XpList) Len() int           { return len(a) }
func (a XpList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a XpList) Less(i, j int) bool { return a[i].Xp > a[j].Xp }

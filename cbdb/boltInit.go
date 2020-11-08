package cbdb

import bolt "go.etcd.io/bbolt"

// BoltInit initialises a bolt database object
func BoltInit(db *bolt.DB) (*BoltDb, error) {
	err := db.Update(func(tx *bolt.Tx) (err error) {
		_, err = tx.CreateBucketIfNotExists([]byte("starboardBlacklist"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("errors"))
		return
	})

	return &BoltDb{Bolt: db}, err
}

// InitForGuild initialises the buckets for a guild
func (db *BoltDb) InitForGuild(guildID string) (err error) {
	err = db.Bolt.Update(func(tx *bolt.Tx) (err error) {
		starboardBucket := tx.Bucket([]byte("starboardBlacklist"))
		_, err = starboardBucket.CreateBucketIfNotExists([]byte(guildID))
		return
	})
	return
}

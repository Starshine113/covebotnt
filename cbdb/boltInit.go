package cbdb

import bolt "go.etcd.io/bbolt"

// BoltInit initialises a bolt database object
func BoltInit(db *bolt.DB) (*BoltDb, error) {
	err := db.Update(func(tx *bolt.Tx) (err error) {
		_, err = tx.CreateBucketIfNotExists([]byte("levels"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("starboardBlacklist"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("levelBlacklist"))
		return err
	})

	return &BoltDb{Bolt: db}, err
}

// InitForGuild initialises the buckets for a guild
func (db *BoltDb) InitForGuild(guildID string) (err error) {
	err = db.Bolt.Update(func(tx *bolt.Tx) (err error) {
		levelBucket := tx.Bucket([]byte("levels"))
		starboardBucket := tx.Bucket([]byte("levels"))
		levelBlacklistBucket := tx.Bucket([]byte("levels"))

		_, err = levelBucket.CreateBucketIfNotExists([]byte(guildID))
		if err != nil {
			return
		}
		_, err = starboardBucket.CreateBucketIfNotExists([]byte(guildID))
		if err != nil {
			return
		}
		_, err = levelBlacklistBucket.CreateBucketIfNotExists([]byte(guildID))
		if err != nil {
			return
		}

		return
	})
	return
}

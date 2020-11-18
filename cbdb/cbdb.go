package cbdb

import (
	"github.com/ReneKroon/ttlcache/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	bolt "go.etcd.io/bbolt"
)

// Db gives access to the database
type Db struct {
	Pool       *pgxpool.Pool
	GuildCache *ttlcache.Cache
}

// BoltDb gives access to the bolt database
type BoltDb struct {
	Bolt *bolt.DB
}

package cbdb

import (
	"github.com/jackc/pgx/v4/pgxpool"
	bolt "go.etcd.io/bbolt"
)

// Db gives access to the database
type Db struct {
	Pool *pgxpool.Pool
}

// BoltDb gives access to the bolt database
type BoltDb struct {
	Bolt *bolt.DB
}

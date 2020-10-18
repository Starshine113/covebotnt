package cbdb

import "github.com/jackc/pgx/v4/pgxpool"

// Db gives access to the database
type Db struct {
	Pool *pgxpool.Pool
}

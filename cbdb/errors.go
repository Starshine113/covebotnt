package cbdb

import bolt "go.etcd.io/bbolt"

// CmdError is the error + its ID
type CmdError struct {
	ErrorID string
	Error   string
}

// AddError adds an error to the database
func (db *BoltDb) AddError(cmdErr CmdError) (err error) {
	return db.Bolt.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte("errors"))
		return b.Put([]byte(cmdErr.ErrorID), []byte(cmdErr.Error))
	})
}

// GetError gets an error by ID
func (db *BoltDb) GetError(id string) (cmdErr CmdError, err error) {
	var cmd []byte
	err = db.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("errors"))
		v := b.Get([]byte(id))
		copy(cmd, v)
		return nil
	})
	if err != nil {
		return
	}
	return CmdError{ErrorID: id, Error: string(cmd)}, nil
}

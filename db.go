// Package rumble is a simple key / document store wrapper for boltdb with an API comparable to mgo.

package rumble

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/boltdb/bolt"
)

// New returns a new rumble.DB at the given path with permissions of 0600 and a timeout of 1 second
func New(path string) (*DB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	return &DB{Bolt: db}, err
}

// DB is an abstraction for BoltDB that provides an API compaprable to mgo.
type DB struct {
	Bolt      *bolt.DB
	Marshal   func(interface{}) ([]byte, error)
	Unmarshal func([]byte, interface{}) error
	NewId     func() []byte
}

// checkDriver replaces nil Marshal, Unmarshal and NewId functions with those from the bson package
func (db *DB) checkDriver() {
	if db.Marshal == nil {
		db.Marshal = bson.Marshal
	}

	if db.Unmarshal == nil {
		db.Unmarshal = bson.Unmarshal
	}

	if db.NewId == nil {
		db.NewId = func() []byte { return []byte(bson.NewObjectId()) }
	}
}

// Bucket returns a new Bucket after executing bolt's CreateBucketIfNotExists. Bucket will panic
// on any error from bolt.
func (db *DB) Bucket(name string) *Bucket {
	if err := db.Bolt.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(name)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	return &Bucket{Name: name, db: db}
}

// DeleteBucket returns the result from bolt's DeleteBucket operation
func (db *DB) DeleteBucket(name string) error {
	return db.Bolt.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(name))
	})
}

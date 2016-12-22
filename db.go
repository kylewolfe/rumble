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
	Bolt          *bolt.DB
	MarshalFunc   func(interface{}) ([]byte, error)
	UnmarshalFunc func([]byte, interface{}) error
	NewKeyFunc    func() []byte
}

// checkDriver replaces nil Marshal, Unmarshal and NewId functions with those from the bson package
func (db *DB) checkDriver() {
	if db.MarshalFunc == nil {
		db.MarshalFunc = bson.Marshal
	}

	if db.UnmarshalFunc == nil {
		db.UnmarshalFunc = bson.Unmarshal
	}

	if db.NewKeyFunc == nil {
		db.NewKeyFunc = func() []byte { return []byte(bson.NewObjectId()) }
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

// Buckets returns a slice of *Bucket that are present in the current DB
func (db *DB) Buckets() []*Bucket {
	var buckets []*Bucket
	if err := db.Bolt.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			buckets = append(buckets, &Bucket{Name: string(name), db: db})
			return nil
		})
	}); err != nil {
		panic(err)
	}
	return buckets
}

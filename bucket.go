package rumble

import "github.com/boltdb/bolt"

// Bucket represents a BoltDB bucket. Methods run from Bucket are thread safe as they wrap thread safe operations exposed by BoltDB.
type Bucket struct {
	Name string
	db   *DB
}

// Iterate returns a new Iterator with the given predicate. If a predicate is given, it will be called against every key:value pair and must return
// true for the record to be returned by the iterator. A nil predicate may be given to return all values from the bucket. This is useful when a quick
// check can be done against the data before it is passed to the DB.UnmarshalerFunc
func (b *Bucket) Iterate(predicate func(k, v []byte) bool) *Iterator {
	// a driver is needed for this operation, ensure one is set
	b.db.checkDriver()
	return newIterator(b, predicate)
}

// Put is a wrapper for BoltDB's Update method, equivalant to an upsert.
func (b *Bucket) Put(v ...interface{}) (err error) {
	// TODO: support batching

	// a dirver is needed for this operation, ensure one is set
	b.db.checkDriver()

	for _, doc := range v {
		key := getKey(doc)
		if key == nil || len(key) == 0 {
			key = b.db.NewKeyFunc()
			setKey(key, doc)
		}

		var d []byte
		if d, err = b.db.MarshalFunc(doc); err != nil {
			return err
		}

		if err := b.db.Bolt.Update(func(tx *bolt.Tx) error {
			txBucket := tx.Bucket([]byte(b.Name))
			return txBucket.Put(key, d)
		}); err != nil {
			return err
		}
	}
	return nil
}

// Get is a wrapper for BoltDb's Get method.
func (b *Bucket) Get(k []byte, v interface{}) (err error) {
	// a driver is needed for this operation, ensure one is set
	b.db.checkDriver()

	var d []byte
	b.db.Bolt.View(func(tx *bolt.Tx) error {
		txBucket := tx.Bucket([]byte(b.Name))
		d = txBucket.Get(k)
		return nil
	})

	err = b.db.UnmarshalFunc(d, v)
	setKey(k, v)
	return err
}

// Delete is a wrarpper for BoltDB's Delete method.
func (b *Bucket) Delete(k []byte) (err error) {
	// a driver is needed for tis operation, ensure one is set
	b.db.checkDriver()

	return b.db.Bolt.Update(func(tx *bolt.Tx) error {
		txBucket := tx.Bucket([]byte(b.Name))
		return txBucket.Delete(k)
	})
}

// Count returns the number of entries in the Bucket
func (b *Bucket) Count() int {
	i := 0
	b.db.Bolt.View(func(tx *bolt.Tx) error {
		txBucket := tx.Bucket([]byte(b.Name))
		txBucket.ForEach(func(_ []byte, _ []byte) error {
			i++
			return nil
		})
		return nil
	})
	return i
}

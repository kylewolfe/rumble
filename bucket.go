package rumble

import "github.com/boltdb/bolt"

// Bucket represents a BoltDB bucket. Methods run from Bucket are thread safe as they wrap thread safe operations exposed by BoltDB.
type Bucket struct {
	Name string
	db   *DB
}

// NewIter returns a new Iter.
func (b *Bucket) NewIter() *Iter {
	// a driver is needed for this operation, ensure one is set
	b.db.checkDriver()

	return newIter(b)
}

// Put is a wrapper for BoltDB's Update method, equivalant to an upsert.
func (b *Bucket) Put(v interface{}) (err error) {
	// TODO: support batching

	// a dirver is needed for this operation, ensure one is set
	b.db.checkDriver()

	key := getKey(v)
	if key == nil || len(key) == 0 {
		key = b.db.NewKey()
		setKey(key, v)
	}

	var d []byte
	if d, err = b.db.Marshal(v); err != nil {
		return err
	}

	return b.db.Bolt.Update(func(tx *bolt.Tx) error {
		txBucket := tx.Bucket([]byte(b.Name))
		return txBucket.Put(key, d)
	})
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

	err = b.db.Unmarshal(d, v)
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

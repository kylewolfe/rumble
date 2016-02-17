package rumble

import "github.com/boltdb/bolt"

type Bucket struct {
	Name string
	db   *DB
}

func (b *Bucket) NewIter() *Iter {
	// a driver is needed for this operation, ensure one is set
	b.db.checkDriver()

	return newIter(b)
}

func (b *Bucket) Put(v interface{}) (err error) {
	// TODO: support batching

	// a dirver is needed for this operation, ensure one is set
	b.db.checkDriver()

	id := getId(v)
	if id == nil || len(id) == 0 {
		id = b.db.NewId()
		setId(id, v)
	}

	var d []byte
	if d, err = b.db.Marshal(v); err != nil {
		return err
	}

	return b.db.Bolt.Update(func(tx *bolt.Tx) error {
		txBucket := tx.Bucket([]byte(b.Name))
		return txBucket.Put(id, d)
	})
}

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
	setId(k, v)
	return err
}

func (b *Bucket) Delete(k []byte) (err error) {
	// a driver is needed for tis operation, ensure one is set
	b.db.checkDriver()

	return b.db.Bolt.Update(func(tx *bolt.Tx) error {
		txBucket := tx.Bucket([]byte(b.Name))
		return txBucket.Delete(k)
	})
}

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

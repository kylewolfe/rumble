package rumble

import "github.com/boltdb/bolt"

// Iterator iterates over a channel of KV and Unmarshals the value
type Iterator struct {
	bucket    *Bucket
	buf       chan KV
	predicate func(k, v []byte) bool
}

// Next Unmarshals the current document into the given interface. Next will return false when the
// Pipeline channel has been closed and the last record read.
func (i *Iterator) Next(v interface{}) bool {
	kv, ok := <-i.buf

	if ok {
		// unmarshal the result, panic on error
		if err := i.bucket.db.UnmarshalFunc(kv.Value, v); err != nil {
			panic(err)
		}

		setKey(kv.Key, v)
		return true
	}

	return false
}

// KV represents a key / value pair within BoltDB
type KV struct {
	Key, Value []byte
}

func newIterator(b *Bucket, predicate func(k, v []byte) bool) *Iterator {
	i := &Iterator{b, make(chan KV, 2), predicate} // TODO: bigger or configurable buffer size?

	go func(b *Bucket, i *Iterator) {
		if err := b.db.Bolt.View(func(tx *bolt.Tx) error {
			// feed the iterator
			txBucket := tx.Bucket([]byte(b.Name))
			txBucket.ForEach(func(k []byte, v []byte) error {
				if i.predicate == nil || i.predicate(k, v) {
					i.buf <- KV{k, v} // TODO: handle timeout?
				}
				return nil
			})

			// close the iterator
			close(i.buf)

			return nil
		}); err != nil {
			panic(err)
		}
	}(b, i)

	return i
}

package rumble

import "github.com/boltdb/bolt"

type Iter struct {
	bucket *Bucket
	buf    chan KV
}

func (i *Iter) Pipeline(ch chan KV) *Iter {
	i.buf = ch
	return i
}

func (i *Iter) Next(v interface{}) bool {
	kv, ok := <-i.buf

	if ok {
		// unmarshal the result, panic on error
		if err := i.bucket.db.Unmarshal(kv.Value, v); err != nil {
			panic(err)
		}

		setId(kv.Key, v)
		return true
	}

	return false
}

type KV struct {
	Key, Value []byte
}

func newIter(b *Bucket) *Iter {
	i := &Iter{b, make(chan KV, 2)} // TODO: bigger or configurable buffer size?

	go func(b *Bucket, i *Iter) {
		if err := b.db.Bolt.View(func(tx *bolt.Tx) error {
			// feed the buffer
			txBucket := tx.Bucket([]byte(b.Name))
			txBucket.ForEach(func(k []byte, v []byte) error {
				i.buf <- KV{k, v} // TODO: handle timeout?
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

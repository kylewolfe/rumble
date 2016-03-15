package rumble

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"

	"gopkg.in/mgo.v2/bson"
)

type BenchValue struct {
	S *S
	B []byte
}

type S struct {
	ID       bson.ObjectId `rumble:"key"`
	Foo, Bar string
	Fizz     int
}

var (
	WithIDs []*BenchValue
	NoIDs   []*S
)

func init() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1000; i++ {
		s := &S{ID: bson.NewObjectId(), Foo: fmt.Sprintf("foo-%v", i), Bar: fmt.Sprintf("bar-%v", i), Fizz: i}
		s2 := &S{Foo: fmt.Sprintf("foo-%v", i), Bar: fmt.Sprintf("bar-%v", i), Fizz: i}
		b, err := bson.Marshal(s)
		if err != nil {
			panic(err)
		}
		WithIDs = append(WithIDs, &BenchValue{s, b})
		NoIDs = append(NoIDs, s2)
	}
}

func BenchmarkBoltDBInsert(b *testing.B) {
	db := prepareTestDB(b)
	defer os.Remove(DBPATH)
	if err := db.Bolt.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("bolt"))
		if err != nil {
			b.Fatal(err)
		}
		for n := 0; n < b.N; n++ {
			i := WithIDs[len(WithIDs)-1]
			err := bucket.Put([]byte(i.S.ID), i.B)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkRumbleDBBSONInsertWithID(b *testing.B) {
	db := prepareTestDB(b)
	defer os.Remove(DBPATH)
	bucket := db.Bucket("rumble")
	for n := 0; n < b.N; n++ {
		i := WithIDs[len(WithIDs)-1]
		if err := bucket.Put(i); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRumbleDBBSONInsertWithoutID(b *testing.B) {
	db := prepareTestDB(b)
	defer os.Remove(DBPATH)

	bucket := db.Bucket("rumble")
	for n := 0; n < b.N; n++ {
		i := NoIDs[len(NoIDs)-1]
		if err := bucket.Put(i); err != nil {
			b.Fatal(err)
		}
	}
}

package rumble

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/boltdb/bolt"
)

type Str struct {
	Foo, Bar string
	Fizz     int
	ID       bson.ObjectId `rumble:"key"`
}

var (
	strID   *Str
	strNoID *Str
	value   []byte
	err     error
)

func init() {
	rand.Seed(time.Now().Unix())
	strID = &Str{"foo", "bar", 1, bson.NewObjectId()}
	strNoID = &Str{Foo: "foo", Bar: "bar", Fizz: 1}
	value, err = bson.Marshal(strID)
	if err != nil {
		panic(err)
	}
}

// TODO: Make test closer to what RumbleDB is actually doing
func benchBoltInsert(b *testing.B, batchSize int) {
	db := prepareTestDB(b, "benchBoltInsert.db")
	defer os.Remove("benchBoltInsert.db")

	var keys [][]byte
	for i := 0; i < batchSize; i++ {
		keys = append(keys, []byte(string(strID.ID)+strconv.Itoa(i)))
	}

	for n := 0; n < b.N; n++ {
		if err := db.Bolt.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte("foo"))
			if err != nil {
				return err
			}

			for _, k := range keys {
				if err = bucket.Put(k, value); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			b.Fatal(err)
		}
	}
}

func benchRumbleInsertWithID(b *testing.B, batchSize int) {
	db := prepareTestDB(b, "benchRumbleInsertWithID.db")
	defer os.Remove("benchRumbleInsertWithID.db")

	bucket := db.Bucket("foo")

	var ins []interface{}
	for i := 0; i < batchSize; i++ {
		ins = append(ins, &Str{strID.Foo, strID.Bar, strID.Fizz, strID.ID})
	}

	for n := 0; n < b.N; n++ {
		if err := bucket.Put(ins...); err != nil {
			b.Fatal(err)
		}
	}

}

func benchRumbleInsertNoID(b *testing.B, batchSize int) {
	db := prepareTestDB(b, "benchRumbleInsertNoID.db")
	defer os.Remove("benchRumbleInsertNoID.db")

	bucket := db.Bucket("foo")

	var ins []interface{}
	for i := 0; i < batchSize; i++ {
		ins = append(ins, &Str{Foo: strID.Foo, Bar: strID.Bar, Fizz: strID.Fizz})
	}

	for n := 0; n < b.N; n++ {
		if err := bucket.Put(ins...); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBoltInsert_1(b *testing.B) {
	benchBoltInsert(b, 1)
}

func BenchmarkRumbleInsertWithID_1(b *testing.B) {
	benchRumbleInsertWithID(b, 1)
}

func BenchmarkRumbleInsertNoID_1(b *testing.B) {
	benchRumbleInsertNoID(b, 1)
}

func BenchmarkBoltInsert_10(b *testing.B) {
	benchBoltInsert(b, 10)
}

func BenchmarkRumbleInsertWithID_10(b *testing.B) {
	benchRumbleInsertWithID(b, 10)
}

func BenchmarkRumbleInsertNoID_10(b *testing.B) {
	benchRumbleInsertNoID(b, 10)
}

func BenchmarkBoltInsert_100(b *testing.B) {
	benchBoltInsert(b, 100)
}

func BenchmarkRumbleInsertWithID_100(b *testing.B) {
	benchRumbleInsertWithID(b, 100)
}

func BenchmarkRumbleInsertNoID_100(b *testing.B) {
	benchRumbleInsertNoID(b, 100)
}

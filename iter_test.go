package rumble

import (
	"os"
	"strconv"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

type iterTest struct {
	predicate func(k, v []byte) bool
	result    []string
}

var iterTests = []iterTest{
	{
		nil,
		[]string{"0", "1", "2", "3"},
	},
	{
		func(k, v []byte) bool {
			if string(k) != "2" {
				return true
			}
			return false
		},
		[]string{"0", "1", "3"},
	},
	{
		func(k, v []byte) bool {
			if string(k) == "foo" {
				return true
			}
			return false
		},
		[]string{},
	},
}

func TestIterator(t *testing.T) {
	db := prepareTestDB(t, "TestIter.db")
	defer func() {
		db.Bolt.Close()
		os.Remove("TestIter.db")
	}()

	bucket := db.Bucket("foo")

	for i := 0; i < 4; i++ {
		bucket.Put(bson.M{
			"_key": []byte(strconv.Itoa(i)),
		})
	}

	for _, it := range iterTests {
		i := bucket.Iterator(it.predicate)

		var results []string

		m := make(map[string]interface{})
		for i.Next(&m) {
			results = append(results, string(m["_key"].([]byte)))
		}

		if len(it.result) != len(results) {
			t.Fatalf("expected %v results, got %v", len(it.result), len(results))
		}

		for k, v := range it.result {
			if results[k] != v {
				t.Fatalf("expected %s, got %s (%#v)", v, results[k], it)
			}
		}
	}
}

package rumble

import (
	"os"
	"testing"
)

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			// TODO: error from here does not provide previous line number
			t.Errorf("function expected to panic but did not")
		}
	}()
	f()
}

func prepareTestDB(bt interface{}, name string) *DB {
	db, err := New(name)
	if err != nil {
		switch t := bt.(type) {
		case *testing.T:
			t.Fatal(err)
		case *testing.B:
			t.Fatal(err)
		default:
			panic(err)
		}
	}
	return db
}

func TestDB(t *testing.T) {
	db := prepareTestDB(t, "TestDB.db")
	defer func() {
		db.Bolt.Close()
		os.Remove("TestDB.db")
	}()

	b := db.Bucket("foo")

	if c := b.Count(); c != 0 {
		t.Fatalf("expected count to be 0, got %v", c)
	}

	if err := b.Put(map[string]interface{}{"fizz": "buzz"}); err != nil {
		t.Fatalf("error inserting map into foo: %s", err)
	}

	if c := b.Count(); c != 1 {
		t.Fatalf("expected count to be 1, got %v", c)
	}

	b2 := db.Bucket("foo")
	if c := b2.Count(); c != 1 {
		t.Fatalf("expected count to be 1, got %v", c)
	}

	if err := db.DeleteBucket("foo"); err != nil {
		t.Fatalf("error deleting bucket: %s", err)
	}

	// ensure delete occured by checking for panic on count
	assertPanic(t, func() { b2.Count() })
}

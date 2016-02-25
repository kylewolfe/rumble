package rumble

import (
	"os"
	"testing"
)

const dbPath = "test.db"

var (
	db  *DB
	err error
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

func TestMain(m *testing.M) {
	var r int
	defer os.Exit(r)

	// create test.db
	db, err = New(dbPath)

	// run tests
	r = m.Run()

	// clean up test.db
	os.Remove(dbPath)
}

func TestDB(t *testing.T) {
	b := db.Bucket("foo")

	if c := b.Count(); c != 0 {
		t.Fatalf("expected count to be 0, got %v", c)
	}

	if err = b.Put(map[string]interface{}{"fizz": "buzz"}); err != nil {
		t.Fatalf("error inserting map into foo: %s", err)
	}

	if c := b.Count(); c != 1 {
		t.Fatalf("expected count to be 1, got %v", c)
	}

	b2 := db.Bucket("foo")
	if c := b2.Count(); c != 1 {
		t.Fatalf("expected count to be 1, got %v", c)
	}

	if err = db.DeleteBucket("foo"); err != nil {
		t.Fatalf("error deleting bucket: %s", err)
	}

	assertPanic(t, func() { b2.Count() })
}

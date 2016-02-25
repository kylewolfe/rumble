package rumble

import "testing"

type customIdString string

func (s customIdString) Foo() string {
	return string(s)
}

type customIdByteSlice []byte

func (s customIdByteSlice) Foo() string {
	return string(s)
}

type M map[string]interface{}

func TestGetId(t *testing.T) {
	// test all combos of body: struct, struct pointer, map, map pointer and id: []byte, alias to []byte, string, alias to string
	expected := []byte("123")

	// struct, []byte
	structByte := &struct {
		Id  []byte `rumble:"id"`
		Foo string
	}{
		expected,
		"bar",
	}

	// struct, alias to []byte
	structByteAlias := &struct {
		Id  customIdByteSlice `rumble:"id"`
		Foo string
	}{
		customIdByteSlice(expected),
		"bar",
	}

	// struct, string
	structString := &struct {
		Id  string `rumble:"id"`
		Foo string
	}{
		string(expected),
		"bar",
	}

	// struct, alias to string
	structStringAlias := &struct {
		Id  customIdString `rumble:"id"`
		Foo string
	}{
		customIdString(expected),
		"bar",
	}

	// test non pointer structs
	if id := getId(*structByte); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, *structByte)
	}
	if id := getId(*structByteAlias); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, *structByteAlias)
	}
	if id := getId(*structString); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, *structString)
	}
	if id := getId(*structStringAlias); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, *structStringAlias)
	}

	// test pointer structs
	if id := getId(structByte); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, structByte)
	}
	if id := getId(structByteAlias); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, structByteAlias)
	}
	if id := getId(structString); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, structString)
	}
	if id := getId(structStringAlias); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, structStringAlias)
	}

	// map, byte
	mapByte := &M{"_id": expected, "foo": "bar"}

	// map, byte alias
	mapByteAlias := &M{"_id": customIdByteSlice(expected), "foo": "bar"}

	// map, string
	mapString := &M{"_id": customIdByteSlice(expected), "foo": "bar"}

	// map, string alias
	mapStringAlias := &M{"_id": customIdByteSlice(expected), "foo": "bar"}

	// test non pointer maps
	if id := getId(*mapByte); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, *mapByte)
	}
	if id := getId(*mapByteAlias); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, *mapByteAlias)
	}
	if id := getId(*mapString); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, *mapString)
	}
	if id := getId(*mapStringAlias); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, *mapStringAlias)
	}

	// test pointer maps
	if id := getId(mapByte); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, mapByte)
	}
	if id := getId(mapByteAlias); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, mapByteAlias)
	}
	if id := getId(mapString); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, mapString)
	}
	if id := getId(mapStringAlias); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, mapStringAlias)
	}
}

func TestSetId(t *testing.T) {
	// test all combos of body: struct, struct pointer, map, map pointer and id: []byte, alias to []byte, string, alias to string
	expected := []byte("123")

	// struct, []byte
	structByte := &struct {
		Id  []byte `rumble:"id"`
		Foo string
	}{
		Foo: "bar",
	}

	// struct, alias to []byte
	structByteAlias := &struct {
		Id  customIdByteSlice `rumble:"id"`
		Foo string
	}{
		Foo: "bar",
	}

	// struct, string
	structString := &struct {
		Id  string `rumble:"id"`
		Foo string
	}{
		Foo: "bar",
	}

	// struct, alias to string
	structStringAlias := &struct {
		Id  customIdString `rumble:"id"`
		Foo string
	}{
		Foo: "bar",
	}

	// test pointer structs
	setId(expected, structByte)
	if id := getId(structByte); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, structByte)
	}
	setId(expected, structByteAlias)
	if id := getId(structByteAlias); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, structByteAlias)
	}
	setId(expected, structString)
	if id := getId(structString); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, structString)
	}
	setId(expected, structStringAlias)
	if id := getId(structStringAlias); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, structStringAlias)
	}

	// map, byte
	mapByte := &M{"foo": "bar"}

	// map, byte alias
	mapByteAlias := &M{"foo": "bar"}

	// map, string
	mapString := &M{"foo": "bar"}

	// map, string alias
	mapStringAlias := &M{"foo": "bar"}

	// test non pointer maps
	setId(expected, *mapByte)
	if id := getId(*mapByte); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, *mapByte)
	}
	setId(expected, *mapByteAlias)
	if id := getId(*mapByteAlias); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, *mapByteAlias)
	}
	setId(expected, *mapString)
	if id := getId(*mapString); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, *mapString)
	}
	setId(expected, *mapStringAlias)
	if id := getId(*mapStringAlias); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, *mapStringAlias)
	}

	// test pointer maps
	setId(expected, mapByte)
	if id := getId(mapByte); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, mapByte)
	}
	setId(expected, mapByteAlias)
	if id := getId(mapByteAlias); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, mapByteAlias)
	}
	setId(expected, mapString)
	if id := getId(mapString); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, mapString)
	}
	setId(expected, mapStringAlias)
	if id := getId(mapStringAlias); string(id) != string(expected) {
		t.Fatalf("expected id '%s', got '%s' from %#v", expected, id, mapStringAlias)
	}
}

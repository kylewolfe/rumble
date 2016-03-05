package rumble

import "testing"

type customKeyString string

func (s customKeyString) Foo() string {
	return string(s)
}

type customKeyByteSlice []byte

func (s customKeyByteSlice) Foo() string {
	return string(s)
}

type M map[string]interface{}

func TestGetKey(t *testing.T) {
	// test all combos of body: struct, struct pointer, map, map pointer and key: []byte, alias to []byte, string, alias to string
	expected := []byte("123")

	// struct, []byte
	structByte := &struct {
		Key  []byte `rumble:"key"`
		Foo string
	}{
		expected,
		"bar",
	}

	// struct, alias to []byte
	structByteAlias := &struct {
		Key  customKeyByteSlice `rumble:"key"`
		Foo string
	}{
		customKeyByteSlice(expected),
		"bar",
	}

	// struct, string
	structString := &struct {
		Key  string `rumble:"key"`
		Foo string
	}{
		string(expected),
		"bar",
	}

	// struct, alias to string
	structStringAlias := &struct {
		Key  customKeyString `rumble:"key"`
		Foo string
	}{
		customKeyString(expected),
		"bar",
	}

	// test non pointer structs
	if key := getKey(*structByte); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, *structByte)
	}
	if key := getKey(*structByteAlias); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, *structByteAlias)
	}
	if key := getKey(*structString); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, *structString)
	}
	if key := getKey(*structStringAlias); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, *structStringAlias)
	}

	// test pointer structs
	if key := getKey(structByte); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, structByte)
	}
	if key := getKey(structByteAlias); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, structByteAlias)
	}
	if key := getKey(structString); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, structString)
	}
	if key := getKey(structStringAlias); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, structStringAlias)
	}

	// map, byte
	mapByte := &M{"_key": expected, "foo": "bar"}

	// map, byte alias
	mapByteAlias := &M{"_key": customKeyByteSlice(expected), "foo": "bar"}

	// map, string
	mapString := &M{"_key": customKeyByteSlice(expected), "foo": "bar"}

	// map, string alias
	mapStringAlias := &M{"_key": customKeyByteSlice(expected), "foo": "bar"}

	// test non pointer maps
	if key := getKey(*mapByte); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, *mapByte)
	}
	if key := getKey(*mapByteAlias); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, *mapByteAlias)
	}
	if key := getKey(*mapString); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, *mapString)
	}
	if key := getKey(*mapStringAlias); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, *mapStringAlias)
	}

	// test pointer maps
	if key := getKey(mapByte); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, mapByte)
	}
	if key := getKey(mapByteAlias); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, mapByteAlias)
	}
	if key := getKey(mapString); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, mapString)
	}
	if key := getKey(mapStringAlias); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, mapStringAlias)
	}
}

func TestSetKey(t *testing.T) {
	// test all combos of body: struct, struct pointer, map, map pointer and key: []byte, alias to []byte, string, alias to string
	expected := []byte("123")

	// struct, []byte
	structByte := &struct {
		Key  []byte `rumble:"key"`
		Foo string
	}{
		Foo: "bar",
	}

	// struct, alias to []byte
	structByteAlias := &struct {
		Key  customKeyByteSlice `rumble:"key"`
		Foo string
	}{
		Foo: "bar",
	}

	// struct, string
	structString := &struct {
		Key  string `rumble:"key"`
		Foo string
	}{
		Foo: "bar",
	}

	// struct, alias to string
	structStringAlias := &struct {
		Key  customKeyString `rumble:"key"`
		Foo string
	}{
		Foo: "bar",
	}

	// test pointer structs
	setKey(expected, structByte)
	if key := getKey(structByte); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, structByte)
	}
	setKey(expected, structByteAlias)
	if key := getKey(structByteAlias); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, structByteAlias)
	}
	setKey(expected, structString)
	if key := getKey(structString); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, structString)
	}
	setKey(expected, structStringAlias)
	if key := getKey(structStringAlias); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, structStringAlias)
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
	setKey(expected, *mapByte)
	if key := getKey(*mapByte); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, *mapByte)
	}
	setKey(expected, *mapByteAlias)
	if key := getKey(*mapByteAlias); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, *mapByteAlias)
	}
	setKey(expected, *mapString)
	if key := getKey(*mapString); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, *mapString)
	}
	setKey(expected, *mapStringAlias)
	if key := getKey(*mapStringAlias); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, *mapStringAlias)
	}

	// test pointer maps
	setKey(expected, mapByte)
	if key := getKey(mapByte); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, mapByte)
	}
	setKey(expected, mapByteAlias)
	if key := getKey(mapByteAlias); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, mapByteAlias)
	}
	setKey(expected, mapString)
	if key := getKey(mapString); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, mapString)
	}
	setKey(expected, mapStringAlias)
	if key := getKey(mapStringAlias); string(key) != string(expected) {
		t.Fatalf("expected key '%s', got '%s' from %#v", expected, key, mapStringAlias)
	}
}

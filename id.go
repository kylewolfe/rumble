package rumble

import (
	"reflect"
	"strings"
)

func getId(v interface{}) []byte {
	val := getReflectValue(reflect.ValueOf(v))

	// return nil on invalid reflect value
	if !val.IsValid() {
		return nil
	}

	switch val.Kind() {
	case reflect.Slice: // TODO: correct way to detect byte slice?
		return val.Bytes()
	case reflect.Interface:
		switch val.Elem().Kind() {
		case reflect.Slice: // TODO: correct way to detect byte slice?
			return val.Elem().Bytes()
		}
	case reflect.String:
		return []byte(val.String())
	}

	return nil
}

func setId(id []byte, v interface{}) {
	// explicitly set _id for maps
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		rv.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(id))
		return
	case reflect.Ptr:
		elem := rv.Elem()
		if elem.Kind() == reflect.Map {
			elem.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(id))
			return
		}
	}

	// try to find a rumble id tag
	val := getReflectValue(rv)

	// return nil on invalid reflect value
	if !val.IsValid() || !val.CanSet() {
		return
	}

	switch val.Kind() {
	case reflect.Slice: // TODO: correct way to detect byte slice?
		val.SetBytes(id)
	case reflect.String:
		val.SetString(string(id))
	}
}

func getReflectValue(rv reflect.Value) reflect.Value {
	switch rv.Kind() {
	case reflect.Map:
		for _, mk := range rv.MapKeys() {
			if mk.Kind() == reflect.String && mk.String() == "_id" {
				return rv.MapIndex(mk)
			}
		}
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			tag := rv.Type().Field(i).Tag.Get("rumble")
			if strings.Contains(strings.ToLower(tag), "id") {
				return rv.Field(i)
			}
		}
	case reflect.Ptr:
		return getReflectValue(rv.Elem())
	}

	// return blank value struct, can use IsValid() on it
	return reflect.Value{}
}

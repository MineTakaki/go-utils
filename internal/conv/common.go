package conv

import "reflect"

// IsEmpty nilまたは空文字かを判定します
func IsEmpty(i interface{}) bool {
	if i == nil {
		return true
	}

	rv := UnwrapNullable(reflect.ValueOf(i))
	if !rv.IsValid() {
		return true
	}

	switch rv.Kind() {
	case reflect.String, reflect.Array, reflect.Slice:
		return rv.Len() == 0
	}
	return false
}

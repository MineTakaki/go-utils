package conv

import (
	"database/sql"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func _toInt64(v reflect.Value) (int64, bool) {
	convFloat := func(f float64) (int64, bool) {
		if f >= math.MinInt64 && f <= math.MaxInt64 {
			return int64(f), true
		}
		return 0, false
	}
	convString := func(s string) (int64, bool) {
		if n, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64); err == nil {
			return n, true
		}
		return 0, false
	}

	switch v.Kind() {
	case reflect.Invalid:
		return 0, false
	case reflect.Ptr:
		if !v.IsNil() {
			return _toInt64(v.Elem())
		}
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		return v.Int(), true
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint8, reflect.Uintptr:
		return int64(v.Uint()), true
	case reflect.Uint64:
		if u := v.Uint(); u <= math.MaxInt64 {
			return int64(u), true
		}
	case reflect.Float32, reflect.Float64:
		return convFloat(v.Float())
	case reflect.String:
		return convString(v.String())
	}

	if v.Type().ConvertibleTo(int64type) {
		return v.Convert(int64type).Int(), true
	}
	if v.Type().ConvertibleTo(nullInt64Type) {
		x, _ := v.Convert(nullInt64Type).Interface().(sql.NullInt64)
		if x.Valid {
			return x.Int64, true
		}
		return 0, false
	}
	if v.Type().ConvertibleTo(nullInt32Type) {
		x, _ := v.Convert(nullInt32Type).Interface().(sql.NullInt32)
		if x.Valid {
			return int64(x.Int32), true
		}
		return 0, false
	}
	if v.Type().ConvertibleTo(nullFloat64Type) {
		x, _ := v.Convert(nullFloat64Type).Interface().(sql.NullFloat64)
		if x.Valid {
			return convFloat(x.Float64)
		}
		return 0, false
	}
	if v.Type().ConvertibleTo(nullStringType) {
		x, _ := v.Convert(nullStringType).Interface().(sql.NullString)
		if x.Valid {
			return convString(x.String)
		}
		return 0, false
	}

	if v.NumMethod() != 0 {
		var b, ok bool
		for _, name := range []string{"Int", "Int64"} {
			v, b, ok = convertMethod(v, name, int64type)
			if ok {
				return v.Int(), b
			}
		}
	}

	return 0, false
}

// Int64 int64型へ変換します
func Int64(i interface{}) (int64, bool) {
	return _toInt64(reflect.ValueOf(i))
}

// Int int型へ変換します
func Int(i interface{}) (int, bool) {
	n, ok := _toInt64(reflect.ValueOf(i))
	if ok {
		if n >= math.MinInt32 && n <= math.MaxInt32 {
			return int(n), ok
		}
	}
	return 0, false
}

package conv

import (
	"database/sql"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/MineTakaki/go-utils/errors"
)

func _toFloat64(v reflect.Value) (float64, bool) {

	convString := func(s string) (float64, bool) {
		if n, err := strconv.ParseFloat(strings.TrimSpace(s), 64); err == nil {
			return n, true
		}
		return 0, false
	}

	switch v.Kind() {
	case reflect.Invalid:
		return 0, false
	case reflect.Ptr:
		if !v.IsNil() {
			return _toFloat64(v.Elem())
		}
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		return float64(v.Int()), true
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uintptr:
		return float64(v.Uint()), true
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	case reflect.String:
		return convString(v.String())
	}

	if v.Type().ConvertibleTo(int64type) {
		return float64(v.Convert(int64type).Int()), true
	}
	if v.Type().ConvertibleTo(nullInt64Type) {
		x, _ := v.Convert(nullInt64Type).Interface().(sql.NullInt64)
		return float64(x.Int64), x.Valid
	}
	if v.Type().ConvertibleTo(nullInt32Type) {
		x, _ := v.Convert(nullInt32Type).Interface().(sql.NullInt32)
		return float64(x.Int32), x.Valid
	}
	if v.Type().ConvertibleTo(nullFloat64Type) {
		x, _ := v.Convert(nullFloat64Type).Interface().(sql.NullFloat64)
		return x.Float64, x.Valid
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
		for _, name := range []string{"Float", "Float64"} {
			v, b, ok = convertMethod(v, name, float64type)
			if ok {
				return v.Float(), b
			}
		}
	}

	return 0, false
}

// ScanFloat64 float64型に変換します
func ScanFloat64(unk interface{}) (float64, error) {
	switch i := unk.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case string:
		return strconv.ParseFloat(i, 64)
	default:
		v := reflect.ValueOf(unk)
		v = reflect.Indirect(v)
		if v.Type().ConvertibleTo(float64type) {
			fv := v.Convert(float64type)
			return fv.Float(), nil
		} else if v.Type().ConvertibleTo(stringType) {
			sv := v.Convert(stringType)
			s := sv.String()
			return strconv.ParseFloat(s, 64)
		} else {
			return math.NaN(), errors.Errorf("Can't convert %v to float64", v.Type())
		}
	}
}

// Float64 float64型に変換します
func Float64(i interface{}) (float64, bool) {
	return _toFloat64(reflect.ValueOf(i))
	/*
		if f, err := ScanFloat64(i); err == nil {
			return f, true
		}
		return 0, false
	*/
}

// NullFloat64 sql.NullFloat64型に変換します
func NullFloat64(i interface{}) (f sql.NullFloat64, ok bool) {
	if IsEmpty(i) {
		ok = true
		return
	}
	f.Float64, f.Valid = Float64(i)
	ok = f.Valid
	return
}

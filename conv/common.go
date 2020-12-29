package conv

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"math"
	"reflect"
	"strings"

	"github.com/MineTakaki/go-utils/types/decimal"
	orgdec "github.com/shopspring/decimal"
)

type kind int

const (
	invalidKind kind = iota
	boolKind
	complexKind
	intKind
	floatKind
	stringKind
	uintKind
	decimalKind
)

var (
	errBadComparisonType = errors.New("invalid type for comparison")
	errBadComparison     = errors.New("incompatible types for comparison")
	errNoComparison      = errors.New("missing argument for comparison")
)

//IsEmpty nilまたは空文字かを判定します
func IsEmpty(i interface{}) bool {
	if i == nil {
		return true
	}
	if n, ok := i.(sql.NullString); ok {
		return !n.Valid || n.String == ""
	}
	if n, ok := i.(sql.NullInt64); ok {
		return !n.Valid
	}
	if n, ok := i.(sql.NullFloat64); ok {
		return !n.Valid
	}
	if n, ok := i.(sql.NullBool); ok {
		return !n.Valid
	}

	if n, ok := i.(*sql.NullString); ok {
		return !n.Valid || n.String == ""
	}
	if n, ok := i.(*sql.NullInt64); ok {
		return !n.Valid
	}
	if n, ok := i.(*sql.NullFloat64); ok {
		return !n.Valid
	}
	if n, ok := i.(*sql.NullBool); ok {
		return !n.Valid
	}

	if s, ok := i.(string); ok {
		return s == ""
	}
	return false
}

var stringType = reflect.TypeOf((*string)(nil)).Elem()
var int64type = reflect.TypeOf((*int64)(nil)).Elem()
var float64type = reflect.TypeOf((*float64)(nil)).Elem()
var bytesType = reflect.TypeOf((*[]byte)(nil)).Elem()
var nullStringType = reflect.TypeOf((*sql.NullString)(nil)).Elem()
var nullInt64Type = reflect.TypeOf((*sql.NullInt64)(nil)).Elem()
var nullInt32Type = reflect.TypeOf((*sql.NullInt32)(nil)).Elem()
var nullFloat64Type = reflect.TypeOf((*sql.NullFloat64)(nil)).Elem()
var decimalType = reflect.TypeOf((*decimal.Decimal)(nil)).Elem()
var nullDecimalType = reflect.TypeOf((*decimal.NullDecimal)(nil)).Elem()
var decimalType2 = reflect.TypeOf((*orgdec.Decimal)(nil)).Elem()
var nullDecimalType2 = reflect.TypeOf((*orgdec.NullDecimal)(nil)).Elem()
var sqlValuerType = reflect.TypeOf((*driver.Valuer)(nil)).Elem()

//UnwrapNullable sql.NullString等のNullableな型を考慮して値を取得します
func UnwrapNullable(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return reflect.Value{}
		}
		return UnwrapNullable(v.Elem())
	case reflect.Struct:
		if v.NumField() < 2 {
			return v
		}
		f2, ok := v.Type().FieldByNameFunc(func(name string) bool {
			return name == "Valid"
		})
		if !ok || f2.Type.Kind() != reflect.Bool {
			return v
		}
		if !v.FieldByIndex(f2.Index).Bool() {
			return reflect.Value{}
		}
		f1 := v.Type().Field(0)
		if f1.Name == f2.Name {
			f1 = v.Type().Field(1)
		}
		return v.FieldByIndex(f1.Index)
	}
	return v
}

func basicKind(v reflect.Value) kind {
	switch v.Kind() {
	case reflect.Bool:
		return boolKind
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intKind
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintKind
	case reflect.Float32, reflect.Float64:
		return floatKind
	case reflect.Complex64, reflect.Complex128:
		return complexKind
	case reflect.String:
		return stringKind
	}

	switch {
	case v.Type().ConvertibleTo(decimalType):
		return decimalKind
	case v.Type().ConvertibleTo(decimalType2):
		return decimalKind
	}
	return invalidKind
}

func toDecimal(r reflect.Value) decimal.Decimal {
	if r.Type().ConvertibleTo(decimalType) {
		d, ok := r.Convert(decimalType).Interface().(decimal.Decimal)
		if !ok {
			panic("cannot convertible decimal.Decimal")
		}
		return d
	}
	d, ok := r.Convert(decimalType2).Interface().(orgdec.Decimal)
	if !ok {
		panic("cannot convertible decimal.Decimal")
	}
	return decimal.Decimal{d}
}

func compareInt(n1, n2 int64) int {
	if n1 > n2 {
		return 1
	}
	if n1 < n2 {
		return -1
	}
	return 0
}

func compareUint(n1, n2 uint64) int {
	if n1 > n2 {
		return 1
	}
	if n1 < n2 {
		return -1
	}
	return 0
}

func compareFloat(n1, n2 float64) int {
	if n1 > n2 {
		return 1
	}
	if n1 < n2 {
		return -1
	}
	return 0
}

//CompareReflectValue reflect.Valueで同士で大小比較します
func CompareReflectValue(r1, r2 reflect.Value) (int, error) {
	r1 = UnwrapNullable(r1)
	r2 = UnwrapNullable(r2)

	k1 := basicKind(r1)
	k2 := basicKind(r2)

	if k1 == k2 {
		switch k1 {
		case boolKind, complexKind:
			return 0, errBadComparisonType
		case intKind:
			return compareInt(r1.Int(), r2.Int()), nil
		case uintKind:
			return compareUint(r1.Uint(), r2.Uint()), nil
		case floatKind:
			return compareFloat(r1.Float(), r2.Float()), nil
		case stringKind:
			return strings.Compare(r1.String(), r2.String()), nil
		case decimalKind:
			return toDecimal(r1).Cmp(toDecimal(r2)), nil
		}
	}

	switch k1 {
	case boolKind, complexKind:
		return 0, errBadComparisonType
	case intKind:
		switch k2 {
		case uintKind:
			u := r2.Uint()
			if u > math.MaxInt64 {
				return -1, nil
			}
			return compareInt(r1.Int(), int64(u)), nil
		case floatKind:
			return compareFloat(float64(r1.Int()), r2.Float()), nil
		case decimalKind:
			return decimal.NewFromInt(r1.Int()).Cmp(toDecimal(r2)), nil
		}
	case uintKind:
		switch k2 {
		case intKind:
			d := r2.Int()
			if d < 0 {
				return 1, nil
			}
			return compareUint(r1.Uint(), uint64(d)), nil
		case floatKind:
			return compareFloat(float64(r1.Uint()), r2.Float()), nil
		case decimalKind:
			return decimal.NewFromFloat(float64(r1.Uint())).Cmp(toDecimal(r2)), nil
		}
	case floatKind:
		switch k2 {
		case intKind:
			return compareFloat(r1.Float(), float64(r2.Int())), nil
		case uintKind:
			return compareFloat(r1.Float(), float64(r2.Uint())), nil
		case decimalKind:
			return decimal.NewFromFloat(r1.Float()).Cmp(toDecimal(r2)), nil
		}
	case decimalKind:
		switch k2 {
		case intKind:
			return toDecimal(r1).Cmp(decimal.NewFromInt(r2.Int())), nil
		case uintKind:
			return toDecimal(r1).Cmp(decimal.NewFromFloat(float64(r2.Uint()))), nil
		case floatKind:
			return toDecimal(r1).Cmp(decimal.NewFromFloat(r2.Float())), nil
		}
	}
	return 0, errBadComparison
}

//EqualReflectValue reflect.Valueで同士で同一比較します
func EqualReflectValue(r1, r2 reflect.Value) (bool, error) {
	r1 = UnwrapNullable(r1)
	r2 = UnwrapNullable(r2)

	k1 := basicKind(r1)
	k2 := basicKind(r2)

	switch k1 {
	default:
		return false, errBadComparison
	case boolKind:
		switch k2 {
		case boolKind:
			return r1.Bool() == r2.Bool(), nil
		}
	case intKind:
		switch k2 {
		case intKind:
			return r1.Int() == r2.Int(), nil
		case uintKind:
			u := r2.Uint()
			if u > math.MaxInt64 {
				return false, nil
			}
			return r1.Int() == int64(u), nil
		case floatKind:
			return float64(r1.Int()) == r2.Float(), nil
		case decimalKind:
			return decimal.NewFromInt(r1.Int()).Equal(toDecimal(r2)), nil
		}
	case uintKind:
		switch k2 {
		case uintKind:
			return r1.Uint() == r2.Uint(), nil
		case intKind:
			d := r2.Int()
			if d < 0 {
				return false, nil
			}
			return r1.Uint() == uint64(d), nil
		case floatKind:
			return float64(r1.Uint()) == r2.Float(), nil
		case decimalKind:
			return decimal.NewFromFloat(float64(r1.Uint())).Equal(toDecimal(r2)), nil
		}
	case floatKind:
		switch k2 {
		case floatKind:
			return r1.Float() == r2.Float(), nil
		case intKind:
			return r1.Float() == float64(r2.Int()), nil
		case uintKind:
			return r1.Float() == float64(r2.Uint()), nil
		case decimalKind:
			return decimal.NewFromFloat(r1.Float()).Equal(toDecimal(r2)), nil
		}
	case decimalKind:
		switch k2 {
		case decimalKind:
			return toDecimal(r1).Equal(toDecimal(r2)), nil
		case intKind:
			return toDecimal(r1).Equal(decimal.NewFromInt(r2.Int())), nil
		case uintKind:
			return toDecimal(r1).Equal(decimal.NewFromFloat(float64(r2.Uint()))), nil
		case floatKind:
			return toDecimal(r1).Equal(decimal.NewFromFloat(r2.Float())), nil
		}
	case complexKind:
		switch k2 {
		case complexKind:
			return r1.Complex() == r2.Complex(), nil
		}
	case stringKind:
		switch k2 {
		case stringKind:
			return r1.String() == r2.String(), nil
		}
	}
	return false, errBadComparison
}

//Compare a と b を比較して結果を返します
func Compare(i1, i2 interface{}) (c int, err error) {
	return CompareReflectValue(reflect.ValueOf(i1), reflect.ValueOf(i2))
}

//Less a < b の時に true を返します
func Less(a, b interface{}) (f bool, err error) {
	var c int
	c, err = Compare(a, b)
	if err != nil {
		return
	}
	f = c < 0
	return
}

//Grater a > b の時に true を返します
func Grater(a, b interface{}) (f bool, err error) {
	var c int
	c, err = Compare(a, b)
	if err != nil {
		return
	}
	f = c > 0
	return
}

//Equal i1 == i2 の時に true を返します
func Equal(i1, i2 interface{}) (f bool, err error) {
	return EqualReflectValue(reflect.ValueOf(i1), reflect.ValueOf(i2))
}

//convertMethod 指定した名前と戻り値の型で一致するメソッドを探して実行します。
func convertMethod(v reflect.Value, name string, t reflect.Type) (reflect.Value, bool, bool) {
	if m := v.MethodByName(name); !m.IsValid() {
	} else if m.Type().NumIn() != 0 {
	} else if n := m.Type().NumOut(); n == 1 {
		ot := m.Type().Out(0)
		if ot.AssignableTo(t) {
			res := m.Call([]reflect.Value{})
			return res[0], true, true
		}
	} else if n == 2 {
		ot := []reflect.Type{m.Type().Out(0), m.Type().Out(1)}
		if ot[0].AssignableTo(t) && ot[1].Kind() == reflect.Bool {
			res := m.Call([]reflect.Value{})
			return res[0], res[1].Bool(), true
		}
	}
	return v, false, false
}

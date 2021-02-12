package decimal

import (
	"math/big"
	"reflect"

	"github.com/MineTakaki/go-utils/internal/conv"
	"github.com/shopspring/decimal"
)

var stringType = reflect.TypeOf((*string)(nil)).Elem()
var bytesType = reflect.TypeOf((*[]byte)(nil)).Elem()
var decimalType = reflect.TypeOf((*Decimal)(nil)).Elem()
var nullDecimalType = reflect.TypeOf((*NullDecimal)(nil)).Elem()
var decimalType2 = reflect.TypeOf((*decimal.Decimal)(nil)).Elem()
var nullDecimalType2 = reflect.TypeOf((*decimal.NullDecimal)(nil)).Elem()

//SumN 集計します。エラーの場合はNullとしてあつかいます
func SumN(arr ...interface{}) (sum NullDecimal) {
	return
}

//MulN NullDecimal同士で乗算します。エラーの場合はNullとしてあつかいます
func MulN(arr ...interface{}) (m NullDecimal) {
	return
}

//MaxN Nullを除く最も大きな値を返します
func MaxN(arr ...interface{}) (max NullDecimal) {
	return
}

//MinN Nullを除く最も小さな値を返します
func MinN(arr ...interface{}) (min NullDecimal) {
	return
}

func unquoteIfQuoted(arr []byte) string {
	if len(arr) >= 2 {
		if arr[0] == '"' && arr[len(arr)-1] == '"' {
			arr = arr[1 : len(arr)-1]
		}
	}
	return string(arr)
}

//ValueOfWithRV Decimal型で返します
func ValueOfWithRV(rv reflect.Value) (Decimal, bool) {
	switch rv.Kind() {
	case reflect.Invalid:
		return Zero, false
	case reflect.Ptr:
		if !rv.IsNil() {
			return ValueOfWithRV(rv.Elem())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewFromInt(rv.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n := new(big.Int)
		n.SetUint64(rv.Uint())
		return NewFromBigInt(n, 0), true
	case reflect.Float32, reflect.Float64:
		return NewFromFloat(rv.Float()), true
	case reflect.String:
		if d, err := NewFromString(unquoteIfQuoted([]byte(rv.String()))); err == nil {
			return d, true
		}
		return Zero, false
	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			if d, err := NewFromString(unquoteIfQuoted(rv.Bytes())); err == nil {
				return d, true
			}
			return Zero, false
		}
	}

	typ := rv.Type()
	switch {
	case typ.ConvertibleTo(bytesType):
		return ValueOfWithRV(rv.Convert(bytesType))
	case typ.ConvertibleTo(stringType):
		return ValueOfWithRV(rv.Convert(stringType))
	case typ.ConvertibleTo(decimalType):
		if x, ok := rv.Convert(decimalType).Interface().(Decimal); ok {
			return x, true
		}
	case typ.ConvertibleTo(nullDecimalType):
		if x, ok := rv.Convert(nullDecimalType).Interface().(NullDecimal); ok {
			return x.Decimal, x.Valid
		}
	case typ.ConvertibleTo(decimalType2):
		if x, ok := rv.Convert(decimalType2).Interface().(decimal.Decimal); ok {
			return Decimal{x}, true
		}
	case typ.ConvertibleTo(nullDecimalType2):
		if x, ok := rv.Convert(nullDecimalType2).Interface().(decimal.NullDecimal); ok {
			return Decimal{x.Decimal}, x.Valid
		}
	}
	return Zero, false
}

//ValueOf Decimal型で返します
func ValueOf(value interface{}) (Decimal, bool) {
	if value == nil {
		return Decimal{}, false
	}
	return ValueOfWithRV(conv.UnwrapNullable(reflect.ValueOf(value)))
}

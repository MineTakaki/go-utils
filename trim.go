package utils

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/MineTakaki/go-utils/errors"
)

// TrimFuncs string項目に対してTrimSpaceを実施します
func TrimFuncs(i interface{}, fn func(rune) bool) (err error) {
	if i == nil || fn == nil {
		return
	}

	var fnProc func(rv reflect.Value)
	fnProc = func(v reflect.Value) {
		switch v.Kind() {
		default:
			err = errors.Errorf("argument must be reference or struct")
			return
		case reflect.Ptr:
			if v.IsNil() {
				return
			}
			x := v.Elem()
			switch x.Kind() {
			case reflect.String:
				p := v.Interface().(*string)
				*p = strings.TrimFunc(*p, fn)
				return
			case reflect.Struct:
				fnProc(x)
			}
		case reflect.Struct:
			v.CanSet()
			for i, num := 0, v.NumField(); i < num; i++ {
				f := v.Field(i)
				if !f.IsValid() {
					continue
				}
				switch f.Kind() {
				case reflect.String:
					if f.CanSet() {
						f.SetString(strings.TrimFunc(f.String(), fn))
					}
				case reflect.Ptr, reflect.Struct, reflect.Array, reflect.Slice:
					fnProc(f)
				}
			}
			return
		case reflect.Array, reflect.Slice:
			for i, num := 0, v.Len(); i < num; i++ {
				fnProc(v.Index(i))
			}
		}
	}

	fnProc(reflect.ValueOf(i))
	return
}

// TrimSpaces string項目に対してTrimSpaceを実施します
func TrimSpaces(i interface{}) error {
	return TrimFuncs(i, unicode.IsSpace)
}

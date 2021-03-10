package utils

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

//TrimFuncs string項目に対してTrimSpaceを実施します
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
			for i, num := 0, v.NumField(); i < num; i++ {
				f := v.Field(i)
				if !f.IsValid() || !f.CanSet() {
					continue
				}
				switch f.Kind() {
				case reflect.String:
					f.SetString(strings.TrimFunc(f.String(), fn))
				case reflect.Ptr, reflect.Struct:
					fnProc(f)
				}
			}
			return
		}
	}

	fnProc(reflect.ValueOf(i))
	return
}

//TrimSpaces string項目に対してTrimSpaceを実施します
func TrimSpaces(i interface{}) error {
	return TrimFuncs(i, unicode.IsSpace)
}

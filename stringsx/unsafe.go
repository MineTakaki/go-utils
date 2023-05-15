package stringsx

import (
	"reflect"
	"unsafe"
)

func ToBytes(s string, fn func(b []byte)) {
	if fn == nil {
		return
	}
	if s == "" {
		fn(nil)
		return
	}
	var b []byte
	(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	(*reflect.SliceHeader)(unsafe.Pointer(&b)).Cap = len(s)
	(*reflect.SliceHeader)(unsafe.Pointer(&b)).Len = len(s)
	fn(b)
}

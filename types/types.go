package types

import (
	"reflect"
)

//Interface Type型からインターフェイスのタイプを取得します
func Interface(typ reflect.Type) (ityp reflect.Type) {
	if typ == nil {
		return
	}
	for typ.Kind() == reflect.Ptr {
		ityp = typ
		typ = typ.Elem()
	}
	if ityp == nil {
		ityp = reflect.New(typ).Type()
	}
	return
}

//AsError エラーか判定します
func AsError(typ reflect.Type) bool {
	m, ok := typ.MethodByName("Error")
	return ok &&
		m.Name == "Error" &&
		m.PkgPath == ""
}

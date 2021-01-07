package conv

import "reflect"

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

package scanner

import (
	"reflect"

	"github.com/MineTakaki/go-utils/conv"
	"github.com/MineTakaki/go-utils/errors"
	"github.com/MineTakaki/go-utils/types"
)

type (
	//Scanner スキャナー
	Scanner interface {
		Scan(i interface{}, cols []string) error
	}

	//Scannable スキャナブル
	Scannable interface {
		Scan(interface{}) error
	}

	//ScanFunc フィールドスキャン関数
	ScanFunc func(v reflect.Value, s string) error

	//ScanFuncFactory フィールドスキャン関数ファクトリー
	ScanFuncFactory func(typ reflect.Type, tag string, options []string) (ScanFunc, error)
)

// AsScannable Scannableインターフェイスを持っているか判定します
func AsScannable(typ reflect.Type) bool {
	if typ == nil {
		return false
	}
	typ = types.Interface(typ)
	m, ok := typ.MethodByName("Scan")
	return ok && _isScanMethod(m.Func.Type(), false)
}

func _isScanMethod(t reflect.Type, f bool) bool {
	return t.Kind() == reflect.Func &&
		t.NumOut() == 1 &&
		types.AsError(t.Out(0)) &&
		((f && t.NumIn() == 1 && t.In(0).Kind() == reflect.Interface) ||
			(!f && t.NumIn() == 2 && t.In(0).Kind() == reflect.Ptr && t.In(1).Kind() == reflect.Interface))
}

// Scan ScannableインターフェイスのScanメソッドを実行します
func Scan(v reflect.Value, i interface{}) (err error) {
	if conv.IsEmpty(i) {
		return nil
	}

	//値をセットする為、アドレス参照を取得します
	if v.Kind() == reflect.Struct {
		v = v.Addr()
	}

	m := v.MethodByName("Scan")
	if m.IsNil() || !_isScanMethod(m.Type(), true) {
		return errors.Errorf("Scannable interface not implement : %v", v)
	}
	out := m.Call([]reflect.Value{reflect.ValueOf(i)})

	o := out[0].Interface()
	if o != nil {
		err = o.(error)
		if err == nil {
			err = errors.Errorf("Unkown error : %v", o)
		}
	}
	return
}

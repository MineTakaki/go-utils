package scanner

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type (
	fieldDefT struct {
		col  int
		idx  int
		scan ScanFunc
	}

	withHeadderT struct {
		typ     reflect.Type
		minCols int
		fields  []*fieldDefT
	}
)

func rawStuctType(typ reflect.Type) (reflect.Type, error) {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errors.Errorf("arg is must struc type : %v", typ)
	}
	return typ, nil
}

func (s *withHeadderT) Scan(i interface{}, cols []string) (err error) {
	if i == nil {

	}
	v := reflect.ValueOf(i)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	typ := v.Type()
	if s.typ != typ {
		return errors.Errorf("type unmatch : %v, %v", s.typ, typ)
	}

	//カラム数が足りない場合
	if n := len(cols); n < s.minCols {
		return errors.Errorf("too few data columns : min=%d, len=%d", s.minCols, n)
	}

	for _, d := range s.fields {
		err = d.scan(v.Field(d.idx), cols[d.col])
		if err != nil {
			return
		}
	}
	return nil
}

//DefaultScanFunc 既定のフィールドスキャン関数を取得します
func DefaultScanFunc(typ reflect.Type) (fn ScanFunc, err error) {
	fnInt := func(v reflect.Value, s string) (err error) {
		s = strings.TrimSpace(s)
		var n int64
		if s != "" {
			var bits int
			switch v.Type().Kind() {
			default:
				return errors.Errorf("unkown type : %v", v.Type())
			case reflect.Int64:
				bits = 64
			case reflect.Int, reflect.Int32:
				bits = 32
			case reflect.Int16:
				bits = 16
			case reflect.Int8:
				bits = 8
			}
			var err error
			n, err = strconv.ParseInt(s, 10, bits)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}
		}
		v.SetInt(n)
		return nil
	}
	fnUint := func(v reflect.Value, s string) (err error) {
		s = strings.TrimSpace(s)
		var n uint64
		if s != "" {
			var bits int
			switch v.Type().Kind() {
			default:
				return errors.Errorf("unkown type : %v", v.Type())
			case reflect.Uint64:
				bits = 64
			case reflect.Uint, reflect.Uint32:
				bits = 32
			case reflect.Uint16:
				bits = 16
			case reflect.Uint8:
				bits = 8
			}
			var err error
			n, err = strconv.ParseUint(s, 10, bits)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}
		}
		v.SetUint(n)
		return nil
	}

	switch typ.Kind() {
	case reflect.String:
		fn = func(v reflect.Value, s string) error {
			v.SetString(strings.TrimSpace(s))
			return nil
		}
	case reflect.Int64, reflect.Int, reflect.Int32, reflect.Int16, reflect.Int8:
		fn = fnInt
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		fn = fnUint
	default:
		if AsScannable(typ) {
			fn = func(v reflect.Value, s string) error {
				return Scan(v, strings.TrimSpace(s))
			}
			return
		}
		err = errors.Errorf("unkown type : %v", typ)
	}

	return
}

//WithHeadder ヘッダーを指定してスキャナーを生成します
func WithHeadder(i interface{}, tag string, headders []string, fact ScanFuncFactory) (Scanner, error) {
	typ, err := rawStuctType(reflect.TypeOf(i))
	if err != nil {
		return nil, err
	}

	scans := map[reflect.Type]ScanFunc{}

	x := withHeadderT{typ: typ}

	for i, m := 0, typ.NumField(); i < m; i++ {
		f := typ.Field(i)
		t := f.Tag.Get(tag)
		if t == "" {
			continue
		}

		for col, h := range headders {
			if h == t {
				scan, _ := scans[f.Type]
				if scan == nil {
					if fact != nil {
						scan, err = fact(f.Type, t, nil)
						if err != nil {
							return nil, err
						}
					}
					if scan == nil {
						scan, err = DefaultScanFunc(f.Type)
					}
					if err != nil {
						return nil, err
					}
				}
				if scan != nil {
					if min := col + 1; x.minCols < min {
						x.minCols = min
					}
					x.fields = append(x.fields, &fieldDefT{col: col, idx: i, scan: scan})
				}
				break
			}
		}
	}

	return &x, nil
}

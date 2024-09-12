package scanner

import (
	goerr "errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/MineTakaki/go-utils/errors"
)

type (
	fieldDefT struct {
		name string
		col  int
		idx  int
		req  bool
		scan ScanFunc
	}

	header struct {
		typ    reflect.Type
		tag    string
		fact   ScanFuncFactory
		fields []*fieldDefT
	}
)

// ErrNotFoundField フィールドがありません
var ErrNotFoundField = goerr.New("field not found")

// ErrUnkownType 型が不明です
var ErrUnkownType = goerr.New("unkown type")

// ErrTooShortFields フィールド数が不足しています
var ErrTooShortFields = goerr.New("too short length of fields")

// ErrScanData スキャンエラー（変換エラー）
var ErrScanData = goerr.New("data scan error")

func rawStuctType(typ reflect.Type) (reflect.Type, error) {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errors.Errorf("arg is must struc type : %v", typ)
	}
	return typ, nil
}

// DefaultScanFunc 既定のフィールドスキャン関数を取得します
func DefaultScanFunc(typ reflect.Type) (fn ScanFunc, err error) {
	fnInt := func(v reflect.Value, s string) (err error) {
		s = strings.TrimSpace(s)
		var n int64
		if s != "" {
			var bits int
			switch v.Type().Kind() {
			default:
				return errors.Wrapf(ErrScanData, "unkown type : %v", v.Type())
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
			if n, err = strconv.ParseInt(s, 10, bits); err != nil {
				err = errors.Wrapf(ErrScanData, err.Error())
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
				return errors.Wrapf(ErrScanData, "unkown type : %v", v.Type())
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
			if n, err = strconv.ParseUint(s, 10, bits); err != nil {
				err = errors.Wrapf(ErrScanData, err.Error())
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
		err = errors.Wrapf(ErrScanData, "unkown type : %v", typ)
	}

	return
}

func parseTag(tag string) (name string, opt bool) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], strings.Contains(tag[idx+1:], "required")
	}
	return tag, false
}

func makeScanFields(typ reflect.Type, tagKey string, headers []string, fact ScanFuncFactory) ([]*fieldDefT, error) {
	headerMap := map[string]int{}
	for i, h := range headers {
		if _, ok := headerMap[h]; ok {
			continue
		}
		headerMap[h] = i
	}

	scans := map[reflect.Type]ScanFunc{}

	var fields []*fieldDefT
	var err error
	for i, m := 0, typ.NumField(); i < m; i++ {
		f := typ.Field(i)

		//TAGを取得します
		tag := f.Tag.Get(tagKey)
		if tag == "" || tag == "-" {
			continue
		}
		name, req := parseTag(tag)

		var fdef *fieldDefT
		col, ok := headerMap[name]
		if ok {
			scan := scans[f.Type]
			if scan == nil {
				if fact != nil {
					if scan, err = fact(f.Type, name, nil); err != nil {
						return nil, err
					}
				}
				if scan == nil {
					if scan, err = DefaultScanFunc(f.Type); err != nil {
						return nil, err
					}
				}
			}
			if scan != nil {
				fdef = &fieldDefT{name: name, col: col, idx: i, req: req, scan: scan}
			}
		}
		if fdef != nil {
			fields = append(fields, fdef)
		} else if req {
			return nil, errors.Wrapf(ErrNotFoundField, "field('%s') not found", name)
		}
	}
	return fields, nil
}

// WithHeader ヘッダーを指定してスキャナーを生成します
//
//	headersの指定がnilの場合は最初の Scan() で与えた cols をヘッダとして扱います
func WithHeader(i interface{}, tag string, headers []string, fact ScanFuncFactory) (Scanner, error) {
	typ, err := rawStuctType(reflect.TypeOf(i))
	if err != nil {
		return nil, err
	}

	x := header{typ: typ, tag: tag}

	if headers != nil {
		if x.fields, err = makeScanFields(typ, tag, headers, fact); err != nil {
			return nil, err
		}
	} else {
		x.fact = fact
	}

	return &x, nil
}

func (s *header) Scan(i interface{}, cols []string) (err error) {
	//ヘッダーが読み込まれていいなかった場合は1行目をヘッダーとして処理します
	if s.fields == nil {
		s.fields, err = makeScanFields(s.typ, s.tag, cols, s.fact)
		return
	}

	v := reflect.ValueOf(i)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if typ := v.Type(); s.typ != typ {
		return errors.Wrapf(ErrUnkownType, "type unmatch : %v, %v", s.typ, typ)
	}

	n := len(cols)
	for _, f := range s.fields {
		//カラム数が足りない場合
		if n <= f.col {
			if f.req {
				return errors.Wrapf(ErrTooShortFields, "have no field data of '%s', column=%d ", f.name, f.col)
			}
			continue
		}
		err = f.scan(v.Field(f.idx), cols[f.col])
		if err != nil {
			return
		}
	}
	return nil
}

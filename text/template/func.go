package template

import (
	"strings"
	"text/template"

	"github.com/MineTakaki/go-utils/conv"
)

func joinFnc(list interface{}, sep string) string {
	if t, ok := list.([]string); ok {
		return strings.Join(t, sep)
	}
	return ""
}

//Join []stringを連結します
func Join(m template.FuncMap) template.FuncMap {
	if m == nil {
		m = template.FuncMap{}
	}
	m["join"] = joinFnc
	return m
}

//CompareOperators 標準の比較演算関数を上書きします
func CompareOperators(m template.FuncMap) template.FuncMap {
	if m == nil {
		m = template.FuncMap{}
	}
	m["eq"] = eq
	m["ne"] = ne
	m["lt"] = lt
	m["le"] = le
	m["gt"] = gt
	m["ge"] = ge
	return m
}

func eq(arg1 interface{}, arg2 ...interface{}) (bool, error) {
	for _, v2 := range arg2 {
		b, err := conv.Equal(arg1, v2)
		if err != nil {
			return false, err
		}
		if b {
			return true, nil
		}
	}
	return false, nil
}

func ne(arg1, arg2 interface{}) (bool, error) {
	// != is the inverse of ==.
	equal, err := eq(arg1, arg2)
	return !equal, err
}

// lt evaluates the comparison a < b.
func lt(arg1, arg2 interface{}) (bool, error) {
	c, err := conv.Compare(arg1, arg2)
	if err != nil {
		return false, err
	}
	return c < 0, nil
}

// le evaluates the comparison <= b.
func le(arg1, arg2 interface{}) (bool, error) {
	c, err := conv.Compare(arg1, arg2)
	if err != nil {
		return false, err
	}
	return c <= 0, nil
}

// gt evaluates the comparison a > b.
func gt(arg1, arg2 interface{}) (bool, error) {
	c, err := conv.Compare(arg1, arg2)
	if err != nil {
		return false, err
	}
	return c > 0, nil
}

// ge evaluates the comparison a >= b.
func ge(arg1, arg2 interface{}) (bool, error) {
	c, err := conv.Compare(arg1, arg2)
	if err != nil {
		return false, err
	}
	return c >= 0, nil
}

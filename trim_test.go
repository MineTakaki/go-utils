package utils

import (
	"testing"
	"unicode"
)

func TestTrimFuncs(t *testing.T) {
	x := " abcde fg "
	if err := TrimFuncs(x, unicode.IsSpace); err == nil {
		t.Error("not refernce type is error")
	}
	if err := TrimFuncs(&x, unicode.IsSpace); err != nil {
		t.Errorf("%+v", err)
	}
	t.Logf("'%s'", x)

	type T struct {
		A string
		b string
		C int
		D *string
		E *T
		F []T
		G []*T
	}

	d := "  dddd  "
	z := T{
		A: " zzzz  ",
	}
	y := T{
		A: "     aaaaa ",
		b: "    bbbb ",
		D: &d,
		E: &z,
		F: []T{
			{A: "   xxxx   "},
			{A: "   yyyy   "},
		},
		G: []*T{
			{A: "  ggggg  "},
		},
	}
	if err := TrimFuncs(&y, unicode.IsSpace); err != nil {
		t.Errorf("%+v", err)
	}
	t.Logf("%+v", y)
	t.Logf("d='%s'", d)
	t.Logf("z:%+v", z)
	t.Logf("y.G[0]:%+v", y.G[0])
}

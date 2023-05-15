package stringsx_test

import (
	"testing"

	"github.com/MineTakaki/go-utils/stringsx"
)

func TestLenW(t *testing.T) {
	for _, x := range []struct {
		s string
		n int
	}{
		{"", 0},
		{"abc", 3},
		{"ＡＢＣ", 6},
	} {
		if n := stringsx.LenW(x.s); n != x.n {
			t.Errorf("LenW(%s) %d != %d", x.s, x.n, n)
		}
	}
}

package types_test

import (
	"testing"

	"github.com/MineTakaki/go-utils/types"
)

func TestHmsString(t *testing.T) {
	for _, x := range []struct {
		v   types.Hms
		exp string
	}{
		{0, "000000"},
		{1, "000001"},
		{120159, "120159"},
	} {
		if act := x.v.String(); x.exp != act {
			t.Errorf("Hms(%d) : exp(%s) != act(%s)", x.v, x.exp, act)
		}
	}
}

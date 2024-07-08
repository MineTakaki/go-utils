package types_test

import (
	"testing"

	"github.com/MineTakaki/go-utils/types"
)

func TestHms(t *testing.T) {
	var hms types.Hms

	text := "234859"
	hour := 23
	minute := 48
	second := 59

	if err := hms.Scan(text); err != nil {
		t.Errorf("%v", err)
		return
	}
	if act, exp := hms.Hour(), hour; act != exp {
		t.Errorf("Houre(%v): %d != %d", hms, act, exp)
	}
	if act, exp := hms.Minute(), minute; act != exp {
		t.Errorf("Minute(%v): %d != %d", hms, act, exp)
	}
	if act, exp := hms.Second(), second; act != exp {
		t.Errorf("Minute(%v): %d != %d", hms, act, exp)
	}
	if h, m, s := hms.Part(); h != hour || m != minute || s != second {
		t.Errorf("%d,%d,%d != %d,%d,%d", h, m, s, hour, minute, second)
	}
	if ok, err := hms.Validate(); err != nil {
		t.Errorf("%v", err)
	} else if !ok {
		t.Error("Hms.Validate() is not err, but result is false")
	}

	if err := hms.Scan("257199"); err != nil {
		//Scanではチェックは行わない
		t.Errorf("%v", err)
		return
	}
}

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

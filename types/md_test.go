package types

import (
	"testing"
	"time"
)

func TestAdjustMonth(t *testing.T) {
	// プラス方向の確認
	for x, y, z := 1, 0, 1; x < 100; x++ {
		if z > 12 {
			z = 1
			y++
		}
		if dy, m := AdjustMonth(0, x); m != z || y != dy {
			t.Errorf("%d => %d, %d: expect=%d, %d", x, m, dy, z, y)
			return
			//} else {
			//t.Logf("%d => %d, %d", x, m, dy)
		}
		z++
	}

	// マイナス方向の確認
	for x, y, z := 12, 0, 12; x > -100; x-- {
		if z < 1 {
			z = 12
			y--
		}
		if dy, m := AdjustMonth(0, x); m != z || y != dy {
			t.Errorf("%d => %d, %d: expect=%d, %d", x, m, dy, z, y)
			return
			//} else {
			//t.Logf("%d => %d, %d", x, m, dy)
		}
		z--
	}
}

func TestAdjustMd(t *testing.T) {
	// プラス方向の確認
	for m := 1; m <= 12; m++ {
		for x := 1; x < 365; x++ {
			tm := time.Date(2021, time.Month(m), 1, 0, 0, 0, 0, time.Local)
			tm = tm.AddDate(0, 0, x-1)
			if dm, d := AdjustMd(m, x); d != tm.Day() || dm != int(tm.Month()) {
				t.Errorf("%d, %d => %d, %d: expect=%d, %d", m, x, d, dm, tm.Day(), tm.Month())
				return
				//} else {
				//	t.Logf("%d, %d => %d, %d", m, x, d, dm)
			}
		}
	}

	// マイナス方向の確認
	for m := 1; m <= 12; m++ {
		for x := 1; x >= -365; x-- {
			tm := time.Date(2022, time.Month(m), 1, 0, 0, 0, 0, time.Local)
			tm = tm.AddDate(0, 0, x-1)
			if dm, d := AdjustMd(m, x); d != tm.Day() || dm != int(tm.Month()) {
				t.Errorf("%d, %d => %d, %d: expect=%d, %d", m, x, d, dm, tm.Day(), tm.Month())
				return
				//} else {
				//	t.Logf("%d, %d => %d, %d", m, x, d, dm)
			}
		}
	}
}

func TestFFormatMd(t *testing.T) {
	for _, x := range []struct {
		md   Md
		sep  string
		exp  string
		expZ string
	}{
		{2331, ":", "23:31", "23:31"},
		{102, ".", "01.02", "1.2"},
		{2202, ":", "22:02", "22:2"},
		{12302, ".", "123.02", "123.2"},
		{123, ":", "01:23", "1:23"},
		{102, ".", "01.02", "1.2"},
		{51, ":", "00:51", "0:51"},
		{3, ".", "00.03", "0.3"},
		{0, ":", "", ""},
	} {
		if act := x.md.FormatMd(x.sep, false); x.exp != act {
			t.Errorf("expect(%s) != actual(%s)", x.exp, act)
		}
		if act := x.md.FormatMd(x.sep, true); x.expZ != act {
			t.Errorf("expect(%s) != actual(%s)", x.expZ, act)
		}
	}
}

func TestMdCompare(t *testing.T) {
	for _, x := range []struct {
		a, b Md
		c    int
	}{
		{0, 0, 0},
		{0, 1, -1},
		{1, 0, 1},
		{1, 1, 0},
		{0401, 0401, 0},
		{0401, 0402, -1},
		{0401, 0331, 1},
	} {
		if cmp := x.a.Compare(x.b); x.c != cmp {
			t.Errorf("expect(%d) != actual(%d)", x.c, cmp)
		}
	}
}

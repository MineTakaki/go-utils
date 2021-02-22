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
		if m, dy := AdjustMonth(x); m != z || y != dy {
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
		if m, dy := AdjustMonth(x); m != z || y != dy {
			t.Errorf("%d => %d, %d: expect=%d, %d", x, m, dy, z, y)
			return
			//} else {
			//t.Logf("%d => %d, %d", x, m, dy)
		}
		z--
	}
}

func TestAdjustDay(t *testing.T) {
	// プラス方向の確認
	for m := 1; m <= 12; m++ {
		for x := 1; x < 365; x++ {
			tm := time.Date(2021, time.Month(m), 1, 0, 0, 0, 0, time.Local)
			tm = tm.AddDate(0, 0, x-1)
			xm := (tm.Year()-2021)*12 + int(tm.Month()) - m
			if d, dm := AdjustDay(m, x); d != tm.Day() || dm != xm {
				t.Errorf("%d, %d => %d, %d: expect=%d, %d", m, x, d, dm, tm.Day(), xm)
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
			xm := (tm.Year()-2022)*12 + int(tm.Month()) - m
			if d, dm := AdjustDay(m, x); d != tm.Day() || dm != xm {
				t.Errorf("%d, %d => %d, %d: expect=%d, %d", m, x, d, dm, tm.Day(), xm)
				return
				//} else {
				//	t.Logf("%d, %d => %d, %d", m, x, d, dm)
			}
		}
	}
}

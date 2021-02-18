package types

import (
	"testing"
	"time"
)

func TestYmdScan(t *testing.T) {

	var tm time.Time
	var ymd Ymd
	if err := ymd.Scan(tm); err != nil {
		t.Errorf("%+v", err)
	} else {
		t.Logf("%v", ymd)
	}

	tm = time.Now()
	if err := ymd.Scan(tm); err != nil {
		t.Errorf("%+v", err)
	} else {
		t.Logf("%v", ymd)
	}

}

func TestYmdBetweenMonth(t *testing.T) {
	type testData struct {
		act    bool
		ymd    Ymd
		m1, m2 int
	}
	for _, d := range []testData{
		{true, 20200401, 4, 4},
		{true, 20200401, 4, 5},
		{true, 20200401, 3, 4},
		{false, 20200401, 3, 3},
		{false, 20200401, 5, 5},
		{true, 20200101, 12, 1},
		{true, 20201201, 12, 1},
		{true, 20200101, 12, 2},
		{false, 20201101, 12, 2},
		{false, 20200301, 12, 2},
		{true, 20200201, 12, 2},
		{true, 20200101, 12, 2},
		{true, 20201201, 12, 2},
		{false, 20201101, 12, 2},
		{false, 0, 12, 2},
	} {
		if exp := d.ymd.BetweenMonth(d.m1, d.m2); exp != d.act {
			t.Errorf("exp:%v, %+v", exp, d)
		}
	}
}

package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestYmd(t *testing.T) {
	var ymd Ymd

	if err := ymd.Scan("20190301"); err != nil {
		t.Errorf("%v", err)
		return
	} else if y := ymd.Year(); y != 2019 {
		t.Errorf("%v: %d != 2019", ymd, y)
		return
	} else if m := ymd.Month(); m != 3 {
		t.Errorf("%d != 3", m)
		return
	} else if d := ymd.Day(); d != 1 {
		t.Errorf("%d != 1", d)
		return
	} else if py, pm, pd := ymd.Part(); py != y {
		t.Errorf("%d != %d", py, y)
		return
	} else if pm != m {
		t.Errorf("%d != %d", pm, m)
		return
	} else if pd != d {
		t.Errorf("%d != %d", pd, d)
		return
	} else if ok, err := ymd.Validate(); err != nil {
		t.Errorf("%v", err)
		return
	} else if !ok {
		t.Error("Ymd.Validate() is not err, but result is false")
		return
	}

	if err := ymd.Scan("20190229"); err != nil {
		//Scanでは年月日チェックは行わない
		t.Errorf("%v", err)
		return
	}

	if ok, err := Ymd(20190229).Validate(); err == nil {
		t.Error("validate 20190229 will ok, but not")
		return
	} else if ok {
		t.Error("validate 20190229 will false, but not")
		return
	}

	if err := ymd.Scan("20200229"); err != nil {
		t.Errorf("scan '20200229' : %v", err)
		return
	}

	if ok, err := Ymd(20200229).Validate(); err != nil {
		t.Errorf("%v", err)
		return
	} else if !ok {
		t.Error("validate 20200229 will true, but not")
		return
	}
}

func TestYmdScan(t *testing.T) {
	type T struct {
		v   interface{}
		exp Ymd
		ok  bool
	}

	tm := time.Now()
	n := tm.Year()*10000 + int(tm.Month())*100 + tm.Day()

	for _, x := range []T{
		{nil, 0, true},
		{"", 0, true},
		{[]byte{}, 0, true},
		{time.Time{}, 10101, true},
		{&time.Time{}, 10101, true},
		{tm, Ymd(n), true},
		{&tm, Ymd(n), true},
		{"29991231", 29991231, true},
		{"1", 1, true},
	} {
		var ymd Ymd
		if err := ymd.Scan(x.v); x.ok && err != nil {
			t.Errorf("in=%v, err=%+v", x.v, err)
		} else if !x.ok && err == nil {
			t.Errorf("in=%v, must be error", x.v)
		} else if !x.ok {
			if ymd != 0 {
				t.Errorf("in=%v, exp=0, act=%v", x.v, ymd)
			}
		} else if ymd != x.exp {
			t.Errorf("in=%v, exp=%v, act=%v", x.v, x.exp, ymd)
		}
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

func TestYmdAdd(t *testing.T) {
	type T struct {
		ymd     Ymd
		y, m, d int
		exp     Ymd
		comment string
	}

	for _, x := range []T{
		{20210301, 0, 0, -1, 20210228, "月初から-1日"},
		{20210228, 0, -1, 0, 20210128, "月-1"},
		{20210301, 0, -1, 30, 20210303, "月初から-1か月、日を足して日が月末越えした場合"},
		{20210301, 0, 0, 375, 20220311, "月初から365+10日"},
	} {
		if ymd := x.ymd.Add(x.y, x.m, x.d); ymd != x.exp {
			t.Errorf("%v.Add(%d,%d,%d), expect=%v, actual=%v: %s", x.ymd, x.y, x.m, x.d, x.exp, ymd, x.comment)
		}
	}
}

func TestDays(t *testing.T) {
	now := YmdNow()
	n := 0
	for ymd := Ymd(10101); ymd <= now; ymd = ymd.Next() {
		if days := ymd.Days(); days != n {
			t.Errorf("days no match %d: %d != %d", ymd, days, n)
			break
		}
		n++
	}
}

func TestYmdUnmarshalJson(t *testing.T) {
	type T struct {
		src string
		exp Ymd
		ok  bool
	}
	for _, x := range []T{
		{"\"19980101\"", 19980101, true},
		{"\"1998/01/01\"", 19980101, true},
		{"19980101", 19980101, true},
		{"null", 0, true},
		{"\"\"", 0, true},
	} {
		var ymd Ymd
		if err := json.Unmarshal([]byte(x.src), &ymd); err != nil {
			if x.ok {
				t.Errorf("%+v", err)
			}
			continue
		}
		if ymd != x.exp {
			t.Errorf("%s: exp=%d, act=%d", x.src, x.exp, ymd)
		}
	}
}

func TestFFormatYmd(t *testing.T) {
	for _, x := range []struct {
		ymd  Ymd
		sep  string
		exp  string
		expZ string
	}{
		{20220102, "/", "2022/01/02", "2022/1/2"},
		{20220102, "-", "2022-01-02", "2022-1-2"},
		{20221031, "/", "2022/10/31", "2022/10/31"},
		{20221031, "-", "2022-10-31", "2022-10-31"},
		{1000102, "/", "0100/01/02", "100/1/2"},
		{100102, "-", "0010-01-02", "10-1-2"},
		{10102, "/", "0001/01/02", "1/1/2"},
		{302, "-", "0000-03-02", "0-3-2"},
		{2, "-", "0000-00-02", "0-0-2"},
		{0, "-", "", ""},
	} {
		if act := x.ymd.FormatYmd(x.sep, false); x.exp != act {
			t.Errorf("expect(%s) != actual(%s)", x.exp, act)
		}
		if act := x.ymd.FormatYmd(x.sep, true); x.expZ != act {
			t.Errorf("expect(%s) != actual(%s)", x.expZ, act)
		}
	}
}

func TestYmdCompare(t *testing.T) {
	for _, x := range []struct {
		a, b Ymd
		c    int
	}{
		{0, 0, 0},
		{0, 1, -1},
		{1, 0, 1},
		{1, 1, 0},
		{20220401, 20220401, 0},
		{20220401, 20220402, -1},
		{20220401, 20220331, 1},
	} {
		if cmp := x.a.Compare(x.b); x.c != cmp {
			t.Errorf("expect(%d) != actual(%d)", x.c, cmp)
		}
	}
}

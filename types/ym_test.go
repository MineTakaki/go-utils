package types

import "testing"

func TestYm(t *testing.T) {
	if exp, act := Ym(201903), Ym(201904).Prev(); exp != act {
		t.Errorf("expect is %d, but actual is %d", exp, act)
		return
	}
	if exp, act := Ym(201905), Ym(201904).Next(); exp != act {
		t.Errorf("expect is %d, but actual is %d", exp, act)
		return
	}
	if exp, act := Ym(201812), Ym(201901).Prev(); exp != act {
		t.Errorf("expect is %d, but actual is %d", exp, act)
		return
	}
	if exp, act := Ym(202001), Ym(201912).Next(); exp != act {
		t.Errorf("expect is %d, but actual is %d", exp, act)
		return
	}
	if expFm, expTo, actFm, actTo := func() (a Ymd, b Ymd, c Ymd, d Ymd) {
		a = Ymd(20191201)
		b = Ymd(20191231)
		c, d = Ym(201912).Term()
		return
	}(); expFm != actFm {
		t.Errorf("expect is %d, but actual is %d", expFm, actFm)
		return
	} else if expTo != actTo {
		t.Errorf("expect is %d, but actual is %d", expTo, actTo)
		return
	}

	var fm, to Ym

	if err := fm.Scan("201912"); err != nil {
		t.Errorf("%v", err)
		return
	} else if fm != 201912 {
		t.Errorf("%v != %d", fm, 201912)
		return
	}
	if err := to.Scan("201903"); err != nil {
		t.Errorf("%v", err)
		return
	} else if to != 201903 {
		t.Errorf("%v != %d", to, 201903)
		return
	}
}

//TestYmBetweenMonth Ym型のBetweenMonth()のテストケース
func TestYmBetweenMonth(t *testing.T) {
	type testData struct {
		act    bool
		ym     Ym
		m1, m2 int
	}
	for _, d := range []testData{
		{true, 202004, 4, 4},
		{true, 202004, 4, 5},
		{true, 202004, 3, 4},
		{false, 202004, 3, 3},
		{false, 202004, 5, 5},
		{true, 202001, 12, 1},
		{true, 202012, 12, 1},
		{true, 202001, 12, 2},
		{false, 202011, 12, 2},
		{false, 202003, 12, 2},
		{true, 202002, 12, 2},
		{true, 202001, 12, 2},
		{true, 202012, 12, 2},
		{false, 202011, 12, 2},
		{false, 0, 12, 2},
	} {
		if exp := d.ym.BetweenMonth(d.m1, d.m2); exp != d.act {
			t.Errorf("exp:%v, %+v", exp, d)
		}
	}
}

func TestMd(t *testing.T) {
	var fm, to Md

	if err := fm.Scan("1201"); err != nil {
		t.Errorf("%v", err)
		return
	}

	if err := to.Scan("0131"); err != nil {
		t.Errorf("%v", err)
		return
	}

	if md := Md(1231); !md.Between(fm, to) {
		t.Errorf("(%v <= %v <= %v) is false", fm, md, to)
		return
	}
}

func TestIsLeapYear(t *testing.T) {
	type T struct {
		y int
		x bool
	}

	for _, x := range []T{
		{y: 2000, x: true},
		{y: 2001, x: false},
		{y: 2002, x: false},
		{y: 2003, x: false},
		{y: 2004, x: true},
		{y: 2020, x: true},
		{y: 2021, x: false},
		{y: 2022, x: false},
		{y: 2023, x: false},
		{y: 2024, x: true},
		{y: 2100, x: false},
		{y: 2200, x: false},
		{y: 2300, x: false},
		{y: 2400, x: true},
	} {
		if a := IsLeapYear(x.y); x.x != a {
			t.Errorf("year:%d, expect=%v, actual=%v", x.y, x.x, a)
		}
	}
}

func TestYmAdd(t *testing.T) {
	type T struct {
		ym      Ym
		y, m    int
		exp     Ym
		comment string
	}

	for _, x := range []T{
		{202103, 1, 0, 202203, "+1年"},
		{202103, -1, 0, 202003, "-1年"},
		{202103, 0, +1, 202104, "+1月"},
		{202102, 0, -1, 202101, "-1月"},
		{202103, 0, +9, 202112, "+9月"},
		{202103, 0, +10, 202201, "+10月で年跨ぎ(1月)"},
		{202103, 0, +11, 202202, "+11月で年跨ぎ(2月)"},
		{202103, 0, -3, 202012, "-3月で年跨ぎ(12月)"},
		{202103, 0, -4, 202011, "-4月で年跨ぎ(11月)"},
	} {
		if ym := x.ym.Add(x.y, x.m); ym != x.exp {
			t.Errorf("%v.Add(%d,%d), expect=%v, actual=%v: %s", x.ym, x.y, x.m, x.exp, ym, x.comment)
		}
	}
}

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

func TestMd(t *testing.T) {
	var fm, to MdT

	if err := fm.Scan("1201"); err != nil {
		t.Errorf("%v", err)
		return
	}

	if err := to.Scan("0131"); err != nil {
		t.Errorf("%v", err)
		return
	}

	if md := MdT(1231); !md.Between(fm, to) {
		t.Errorf("(%v <= %v <= %v) is false", fm, md, to)
		return
	}
}

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

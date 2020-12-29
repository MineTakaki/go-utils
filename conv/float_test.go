package conv

import (
	"database/sql"
	"testing"
)

type (
	testFloat struct {
		n float64
	}

	testFloatB struct {
		n float64
		b bool
	}
)

func (t testFloat) Float() float64 {
	return t.n
}

func (t testFloatB) Float() (float64, bool) {
	return t.n, t.b
}

func TestFloat64(t *testing.T) {
	if n, ok := Float64(int64(123)); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Float64(int(123)); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Float64(int8(123)); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Float64(int16(123)); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Float64(int32(123)); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if _, ok := Float64(nil); ok {
		t.Error("error")
	}
	if _, ok := Float64(sql.NullInt64{Int64: 123, Valid: true}); !ok {
		t.Error("error")
	}
	if _, ok := Float64(sql.NullInt64{Int64: 123, Valid: false}); ok {
		t.Error("error")
	}
	if _, ok := Float64(sql.NullFloat64{Float64: 123, Valid: true}); !ok {
		t.Error("error")
	}
	if _, ok := Float64(sql.NullFloat64{Float64: 123, Valid: false}); ok {
		t.Error("error")
	}
	if n, ok := Float64(testFloat{n: 123}); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Float64(&testFloat{n: 123}); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Float64(testFloatB{n: 123, b: true}); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Float64("123.456"); !ok {
		t.Error("error")
	} else if n != 123.456 {
		t.Error("error")
	}
	//Uint
	if n, ok := Float64(uint64(123)); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Float64(uint(123)); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Float64(uint8(123)); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Float64(uint16(123)); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Float64(uint32(123)); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
}

func TestNullFloat64(t *testing.T) {
	if n, ok := NullFloat64("1"); !ok {
		t.Errorf("error")
	} else if !n.Valid {
		t.Errorf("error")
	} else if n.Float64 != 1 {
		t.Errorf("error")
	}
}

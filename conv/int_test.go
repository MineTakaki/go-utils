package conv

import (
	"database/sql"
	"encoding/json"
	"testing"
)

type (
	testInt struct {
		n int64
	}

	testIntB struct {
		n int64
		b bool
	}
)

func (t testInt) Int() int64 {
	return t.n
}

func (t testIntB) Int() (int64, bool) {
	return t.n, t.b
}

func TestInt64(t *testing.T) {
	if n, ok := Int64(int64(123)); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Int64(int(123)); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Int64(int8(123)); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Int64(int16(123)); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Int64(int32(123)); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if _, ok := Int64(nil); ok {
		t.Error("error")
	}
	if _, ok := Int64(sql.NullInt64{Int64: 123, Valid: true}); !ok {
		t.Error("error")
	}
	if _, ok := Int64(sql.NullInt64{Int64: 123, Valid: false}); ok {
		t.Error("error")
	}
	if n, ok := Int64(testInt{n: 123}); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Int64(&testInt{n: 123}); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Int64(testIntB{n: 123, b: true}); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if n, ok := Int64(json.Number("123")); !ok {
		t.Error("error")
	} else if n != 123 {
		t.Error("error")
	}
	if _, ok := Int64(json.Number("")); ok {
		t.Error("error")
	}
}

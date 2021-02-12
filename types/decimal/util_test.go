package decimal

import (
	"database/sql"
	"testing"
)

func TestValueOf(t *testing.T) {
	if d, ok := ValueOf(Cent); !ok {
		t.Error("[Decimal] cannot convert error")
	} else if !d.Equal(Cent) {
		t.Errorf("[Decimal] value unmatch: %v", d)
	}
	if d, ok := ValueOf(100); !ok {
		t.Error("[int] cannot convert error")
	} else if !d.Equal(Cent) {
		t.Errorf("[int] value unmatch: %+v != %+v", d, Cent)
	}
	if d, ok := ValueOf(Null); ok {
		t.Error("[NullDecimal] convert error")
	} else if !d.IsZero() {
		t.Errorf("[Null] value must be Zero: %v", d)
	}
	if d, ok := ValueOf(sql.NullFloat64{}); ok {
		t.Error("[NullFloat64] convert error")
	} else if !d.IsZero() {
		t.Errorf("[Null] value must be Zero: %v", d)
	}
	if d, ok := ValueOf(sql.NullInt32{}); ok {
		t.Error("[NullInt32] convert error")
	} else if !d.IsZero() {
		t.Errorf("[Null] value must be Zero: %v", d)
	}
	if d, ok := ValueOf(sql.NullInt64{}); ok {
		t.Error("[NullInt64] convert error")
	} else if !d.IsZero() {
		t.Errorf("[Null] value must be Zero: %v", d)
	}
	if d, ok := ValueOf(sql.NullString{}); ok {
		t.Error("[NullString] convert error")
	} else if !d.IsZero() {
		t.Errorf("[Null] value must be Zero: %v", d)
	}
}

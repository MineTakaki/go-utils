package decimal_test

import (
	"database/sql"
	"testing"

	"github.com/MineTakaki/go-utils/types/decimal"
)

func TestValueOf(t *testing.T) {
	if d, ok := decimal.ValueOf(decimal.Cent); !ok {
		t.Error("[Decimal] cannot convert error")
	} else if !d.Equal(decimal.Cent) {
		t.Errorf("[Decimal] value unmatch: %v", d)
	}
	if d, ok := decimal.ValueOf(100); !ok {
		t.Error("[int] cannot convert error")
	} else if !d.Equal(decimal.Cent) {
		t.Errorf("[int] value unmatch: %+v != %+v", d, decimal.Cent)
	}
	if d, ok := decimal.ValueOf(decimal.Null); ok {
		t.Error("[NullDecimal] convert error")
	} else if !d.IsZero() {
		t.Errorf("[Null] value must be Zero: %v", d)
	}
	if d, ok := decimal.ValueOf(sql.NullFloat64{}); ok {
		t.Error("[NullFloat64] convert error")
	} else if !d.IsZero() {
		t.Errorf("[Null] value must be Zero: %v", d)
	}
	if d, ok := decimal.ValueOf(sql.NullInt32{}); ok {
		t.Error("[NullInt32] convert error")
	} else if !d.IsZero() {
		t.Errorf("[Null] value must be Zero: %v", d)
	}
	if d, ok := decimal.ValueOf(sql.NullInt64{}); ok {
		t.Error("[NullInt64] convert error")
	} else if !d.IsZero() {
		t.Errorf("[Null] value must be Zero: %v", d)
	}
	if d, ok := decimal.ValueOf(sql.NullString{}); ok {
		t.Error("[NullString] convert error")
	} else if !d.IsZero() {
		t.Errorf("[Null] value must be Zero: %v", d)
	}
}

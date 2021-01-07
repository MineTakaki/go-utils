package decimal

import "testing"

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
}

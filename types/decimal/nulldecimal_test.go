package decimal_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/MineTakaki/go-utils/types/decimal"
)

func TestNullDecimlal_UnmarshallJson(t *testing.T) {
	var doc struct {
		Amount decimal.NullDecimal `json:"amount"`
	}
	for _, docStr := range []string{`{"amount": null}`, `{"amount": ""}`} {
		err := json.Unmarshal([]byte(docStr), &doc)
		if err != nil {
			t.Errorf("error unmarshaling %s: %v", docStr, err)
		} else if doc.Amount.Valid {
			t.Errorf("expected Null, got %s", doc.Amount.String())
		}
	}
}

func TestNullDecimal_Equal(t *testing.T) {
	for _, x := range []struct {
		a, b decimal.NullDecimal
		e    bool
	}{
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			e: true,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(1), Valid: true},
			e: false,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: false},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			e: false,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: false},
			e: false,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(1), Valid: false},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(2), Valid: false},
			e: true,
		},
	} {
		t.Run(
			fmt.Sprintf("{%s,%v}=={%s,%v}", x.a.Decimal.String(), x.a.Valid, x.b.Decimal.String(), x.b.Valid),
			func(t *testing.T) {
				if act := x.a.Equal(x.b); x.e != act {
					t.Errorf("exp(%v) != act(%v)", x.e, act)
				}
				if act := x.b.Equal(x.a); x.e != act {
					t.Errorf("exp(%v) != act(%v)", x.e, act)
				}
			},
		)
	}
}

func TestNullDecimal_EqualNZ(t *testing.T) {
	for _, x := range []struct {
		a, b decimal.NullDecimal
		e    bool
	}{
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			e: true,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(1), Valid: true},
			e: false,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(1), Valid: false},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			e: true,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(1), Valid: false},
			e: true,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(1), Valid: false},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(2), Valid: false},
			e: true,
		},
	} {
		t.Run(
			fmt.Sprintf("{%s,%v}=={%s,%v}", x.a.Decimal.String(), x.a.Valid, x.b.Decimal.String(), x.b.Valid),
			func(t *testing.T) {
				if act := x.a.EqualNZ(x.b); x.e != act {
					t.Errorf("exp(%v) != act(%v)", x.e, act)
				}
				if act := x.b.EqualNZ(x.a); x.e != act {
					t.Errorf("exp(%v) != act(%v)", x.e, act)
				}
			},
		)
	}
}

func TestNullDecimal_Cmp(t *testing.T) {
	for _, x := range []struct {
		a, b decimal.NullDecimal
		e    int
	}{
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			e: 0,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(1), Valid: true},
			e: -1,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(1), Valid: true},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			e: 1,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(1), Valid: false},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			e: -1,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(1), Valid: false},
			e: 1,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(1), Valid: false},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(2), Valid: false},
			e: 0,
		},
	} {
		t.Run(
			fmt.Sprintf("Cmp({%s,%v},{%s,%v})", x.a.Decimal.String(), x.a.Valid, x.b.Decimal.String(), x.b.Valid),
			func(t *testing.T) {
				if act := x.a.Cmp(x.b); x.e != act {
					t.Errorf("exp(%v) != act(%v)", x.e, act)
				}
			},
		)
	}
}

func TestNullDecimal_CmpNZ(t *testing.T) {
	for _, x := range []struct {
		a, b decimal.NullDecimal
		e    int
	}{
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			e: 0,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(1), Valid: true},
			e: -1,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(1), Valid: true},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			e: 1,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(1), Valid: false},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			e: 0,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(0), Valid: true},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(1), Valid: false},
			e: 0,
		},
		{
			a: decimal.NullDecimal{Decimal: decimal.NewFromInt(1), Valid: false},
			b: decimal.NullDecimal{Decimal: decimal.NewFromInt(2), Valid: false},
			e: 0,
		},
	} {
		t.Run(
			fmt.Sprintf("CmpNz({%s,%v},{%s,%v})", x.a.Decimal.String(), x.a.Valid, x.b.Decimal.String(), x.b.Valid),
			func(t *testing.T) {
				if act := x.a.CmpNz(x.b); x.e != act {
					t.Errorf("exp(%v) != act(%v)", x.e, act)
				}
			},
		)
	}
}

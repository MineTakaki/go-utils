package scanner

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestAsScannable(t *testing.T) {
	str := sql.NullString{}
	if !AsScannable(reflect.TypeOf(&str)) {
		t.Error("sql.NullSring is scannable")
	}

	if !AsScannable(reflect.TypeOf(str)) {
		t.Error("sql.NullSring is scannable")
	}

	dec := decimal.Decimal{}
	if !AsScannable(reflect.TypeOf(&dec)) {
		t.Error("decimal.Decimal is scannable")
	}

	if !AsScannable(reflect.TypeOf(dec)) {
		t.Error("decimal.Decimal is scannable")
	}

}

func TestScan(t *testing.T) {
	{
		str := sql.NullString{}
		if err := Scan(reflect.ValueOf(&str), "ABC"); err != nil {
			t.Errorf("%+v", err)
		} else if str.String != "ABC" {
			t.Errorf("%+v", str)
		} else if str.Valid == false {
			t.Errorf("%+v", str)
		} else {
			t.Logf("%+v", str)
		}
		fmt.Printf("%+v\n", str)
	}
	{
		n := sql.NullInt64{}
		if err := Scan(reflect.ValueOf(&n), "123"); err != nil {
			t.Errorf("%+v", err)
		} else if n.Int64 != 123 {
			t.Errorf("%+v", n)
		} else if n.Valid == false {
			t.Errorf("%+v", n)
		} else {
			t.Logf("%+v", n)
		}
		fmt.Printf("%+v\n", n)
	}
	{
		n := sql.NullInt64{}
		if err := Scan(reflect.ValueOf(&n), "0"); err != nil {
			t.Errorf("%+v", err)
		} else if n.Int64 != 0 {
			t.Errorf("%+v", n)
		} else if n.Valid == false {
			t.Errorf("%+v", n)
		} else {
			t.Logf("%+v", n)
		}
		fmt.Printf("%+v\n", n)
	}
	{
		n := sql.NullFloat64{}
		if err := Scan(reflect.ValueOf(&n), "123.456"); err != nil {
			t.Errorf("%+v", err)
		} else if n.Float64 != 123.456 {
			t.Errorf("%+v", n)
		} else if n.Valid == false {
			t.Errorf("%+v", n)
		} else {
			t.Logf("%+v", n)
		}
		fmt.Printf("%+v\n", n)
	}
}

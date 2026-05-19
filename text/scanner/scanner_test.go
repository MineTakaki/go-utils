package scanner_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/MineTakaki/go-utils/text/scanner"
	"github.com/MineTakaki/go-utils/types"
	"github.com/MineTakaki/go-utils/types/decimal"
)

func TestAsScannable(t *testing.T) {
	t.Run("sql.NullString", func(t *testing.T) {
		str := sql.NullString{}
		if !scanner.AsScannable(reflect.TypeOf(&str)) {
			t.Error("sql.NullSring is scannable")
		}

		if !scanner.AsScannable(reflect.TypeOf(str)) {
			t.Error("sql.NullSring is scannable")
		}
	})

	t.Run("decimal.Decimal", func(t *testing.T) {
		dec := decimal.Decimal{}
		if !scanner.AsScannable(reflect.TypeOf(&dec)) {
			t.Error("decimal.Decimal is scannable")
		}

		if !scanner.AsScannable(reflect.TypeOf(dec)) {
			t.Error("decimal.Decimal is scannable")
		}
	})

	t.Run("decimal.NullDecimal", func(t *testing.T) {
		dec := decimal.NullDecimal{}
		if !scanner.AsScannable(reflect.TypeOf(&dec)) {
			t.Error("decimal.NullDecimal is scannable")
		}

		if !scanner.AsScannable(reflect.TypeOf(dec)) {
			t.Error("decimal.NullDecimal is scannable")
		}
	})

	t.Run("types.Ymd", func(t *testing.T) {
		var x types.Ymd
		if !scanner.AsScannable(reflect.TypeOf(&x)) {
			t.Error("types.Ymd is scannable")
		}

		if !scanner.AsScannable(reflect.TypeOf(x)) {
			t.Error("types.Ymd is scannable")
		}
	})

	t.Run("types.Ym", func(t *testing.T) {
		var x types.Ym
		if !scanner.AsScannable(reflect.TypeOf(&x)) {
			t.Error("types.Ym is scannable")
		}

		if !scanner.AsScannable(reflect.TypeOf(x)) {
			t.Error("types.Ym is scannable")
		}
	})

	t.Run("string", func(t *testing.T) {
		var x string
		if scanner.AsScannable(reflect.TypeOf(&x)) {
			t.Error("string is not scannable")
		}
		if scanner.AsScannable(reflect.TypeOf(x)) {
			t.Error("string is not scannable")
		}

	})

	t.Run("int", func(t *testing.T) {
		var x int
		if scanner.AsScannable(reflect.TypeOf(&x)) {
			t.Error("int is not scannable")
		}
		if scanner.AsScannable(reflect.TypeOf(x)) {
			t.Error("int is not scannable")
		}

	})
}

func TestScan(t *testing.T) {
	{
		str := sql.NullString{}
		if err := scanner.Scan(reflect.ValueOf(&str), "ABC"); err != nil {
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
		if err := scanner.Scan(reflect.ValueOf(&n), "123"); err != nil {
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
		if err := scanner.Scan(reflect.ValueOf(&n), "0"); err != nil {
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
		if err := scanner.Scan(reflect.ValueOf(&n), "123.456"); err != nil {
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

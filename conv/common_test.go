package conv

import (
	"database/sql"
	"reflect"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	if !IsEmpty(nil) {
		t.Error("IsEmpty(nil) != true")
	}
	if !IsEmpty("") {
		t.Error("IsEmpty(\"\") != true")
	}
	if IsEmpty(" ") {
		t.Error("IsEmpty(\" \") != false")
	}

	fn := func(exact bool, x interface{}) {
		if exact != IsEmpty(x) {
			t.Errorf("IsEmpty(%+v) != %v", x, exact)
		}
	}

	fn(true, sql.NullString{String: "", Valid: false})
	fn(true, sql.NullString{String: "", Valid: true})
	fn(true, sql.NullString{String: "AAA", Valid: false})
	fn(false, sql.NullString{String: "AAA", Valid: true})
	fn(true, &sql.NullString{String: "", Valid: false})
	fn(true, &sql.NullString{String: "", Valid: true})
	fn(true, &sql.NullString{String: "AAA", Valid: false})
	fn(false, &sql.NullString{String: "AAA", Valid: true})

	fn(true, sql.NullInt64{Int64: 0, Valid: false})
	fn(false, sql.NullInt64{Int64: 0, Valid: true})
	fn(true, sql.NullInt64{Int64: 1, Valid: false})
	fn(false, sql.NullInt64{Int64: 1, Valid: true})
	fn(true, &sql.NullInt64{Int64: 0, Valid: false})
	fn(false, &sql.NullInt64{Int64: 0, Valid: true})
	fn(true, &sql.NullInt64{Int64: 1, Valid: false})
	fn(false, &sql.NullInt64{Int64: 1, Valid: true})

	fn(true, sql.NullFloat64{Float64: 0, Valid: false})
	fn(false, sql.NullFloat64{Float64: 0, Valid: true})
	fn(true, sql.NullFloat64{Float64: 1, Valid: false})
	fn(false, sql.NullFloat64{Float64: 1, Valid: true})
	fn(true, &sql.NullFloat64{Float64: 0, Valid: false})
	fn(false, &sql.NullFloat64{Float64: 0, Valid: true})
	fn(true, &sql.NullFloat64{Float64: 1, Valid: false})
	fn(false, &sql.NullFloat64{Float64: 1, Valid: true})

	fn(true, sql.NullBool{Bool: false, Valid: false})
	fn(false, sql.NullBool{Bool: false, Valid: true})
	fn(true, sql.NullBool{Bool: true, Valid: false})
	fn(false, sql.NullBool{Bool: true, Valid: true})
	fn(true, &sql.NullBool{Bool: false, Valid: false})
	fn(false, &sql.NullBool{Bool: false, Valid: true})
	fn(true, &sql.NullBool{Bool: true, Valid: false})
	fn(false, &sql.NullBool{Bool: true, Valid: true})
}

func TestUnwrapNullable(t *testing.T) {
	if v := UnwrapNullable(reflect.ValueOf(sql.NullInt64{Int64: 123})); v.IsValid() {
		t.Errorf("%+v", v)
	} else if v.IsValid() {
		t.Errorf("%+v", v.Kind())
	}
	if v := UnwrapNullable(reflect.ValueOf(sql.NullInt64{Int64: 123, Valid: true})); !v.IsValid() {
		t.Errorf("%+v", v)
	} else if v.Kind() != reflect.Int64 {
		t.Errorf("%+v", v)
	} else if v.Int() != 123 {
		t.Errorf("%+v", v)
	}
	if v := UnwrapNullable(reflect.ValueOf(&sql.NullInt64{Int64: 123, Valid: true})); !v.IsValid() {
		t.Errorf("%+v", v)
	} else if v.Kind() != reflect.Int64 {
		t.Errorf("%+v", v)
	} else if v.Int() != 123 {
		t.Errorf("%+v", v)
	}
	if v := UnwrapNullable(reflect.ValueOf(sql.NullBool{Bool: true, Valid: true})); !v.IsValid() {
		t.Errorf("%+v", v)
	} else if v.Kind() != reflect.Bool {
		t.Errorf("%+v", v)
	} else if !v.Bool() {
		t.Errorf("%+v", v)
	}
}

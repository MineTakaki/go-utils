package scanner_test

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/MineTakaki/go-utils/text/scanner"
	"github.com/MineTakaki/go-utils/types"
)

func TestWithHeadder(t *testing.T) {
	type testT1 struct {
		Col1 string         `header:"col1"`
		Col2 string         `header:"col2"`
		Col3 string         `json:"hoge" header:"col3"`
		Col4 string         `header:"col4"`
		Col5 sql.NullString `header:"col5"`
		Col6 int            `header:"col6,required"`
		Col7 types.Ymd      `header:"col7"`
	}

	headers := []string{
		"col2",
		"col4",
		"col3",
		"col1",
		"col5",
		"col6",
		"col7",
	}

	rec := testT1{}

	scan, err := scanner.WithHeader(&rec, "header", headers, nil)
	if err != nil {
		t.Errorf("%+v", err)
	}

	err = scan.Scan(&rec, []string{"A", "B", "C", "D", "hhhh", "1234", "2026/12/31"})
	if err != nil {
		t.Errorf("%+v", err)
	}

	fmt.Printf("%+v\n", rec)
}

func TestWithHeadder2(t *testing.T) {
	type testT1 struct {
		Col1 string         `header:"col1"`
		Col2 string         `header:"col2"`
		Col3 string         `header:"[C|c]ol3,regexp,eod"`
		Col4 string         `header:"col4,skip"`
		Col5 sql.NullString `header:"col5"`
		Col6 int            `header:"col6,required,eod"`
	}

	headers := []string{
		"col2",
		"col4",
		"Col3",
		"col1",
		"col5",
		"col6",
	}

	var err error
	var scan scanner.Scanner
	{
		rec := testT1{}
		scan, err = scanner.WithHeader(&rec, "header", headers, nil)
		if err != nil {
			t.Errorf("%+v", err)
		}
		fmt.Printf("%+v\n", rec)
	}

	{
		rec := testT1{}
		err = scan.Scan(&rec, []string{"A", "B", "C", "D", "hhhh", "1234"})
		if err != nil {
			t.Errorf("%+v", err)
		}
		fmt.Printf("%+v\n", rec)
	}

	{
		rec := testT1{}
		err = scan.Scan(&rec, []string{"A", "", "C", "D", "hhhh", "1234"})
		if err == nil {
			t.Error("err is nil")
		} else if !errors.Is(err, scanner.ErrSkipRow) {
			t.Errorf("%+v", err)
		}
		fmt.Printf("%+v\n", rec)
	}

	{
		rec := testT1{}
		err = scan.Scan(&rec, []string{"A", "B", "", "D", "hhhh", ""})
		if err == nil {
			t.Error("err is nil")
		} else if !errors.Is(err, io.EOF) {
			t.Errorf("%+v", err)
		}
		fmt.Printf("%+v\n", rec)
	}

	{
		rec := testT1{}
		err = scan.Scan(&rec, []string{"A", "", "", "D", "hhhh", "1234"})
		if err == nil {
			t.Error("err is nil")
		} else if !errors.Is(err, scanner.ErrSkipRow) {
			t.Errorf("%+v", err)
		}
		fmt.Printf("%+v\n", rec)
	}

	{
		rec := testT1{}
		err = scan.Scan(&rec, []string{"A", "B", "", "D", "hhhh", "1234"})
		if err != nil {
			t.Errorf("%+v", err)
		}
		fmt.Printf("%+v\n", rec)
	}
}

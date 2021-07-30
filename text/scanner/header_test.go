package scanner

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestWithHeadder(t *testing.T) {
	type testT1 struct {
		Col1 string         `header:"col1"`
		Col2 string         `header:"col2"`
		Col3 string         `json:"hoge" header:"col3"`
		Col4 string         `header:"col4"`
		Col5 sql.NullString `header:"col5"`
		Col6 int            `header:"col6,required"`
	}

	headers := []string{
		"col2",
		"col4",
		"col3",
		"col1",
		"col5",
		"col6",
	}

	rec := testT1{}

	scan, err := WithHeader(&rec, "header", headers, nil)
	if err != nil {
		t.Errorf("%+v", err)
	}

	err = scan.Scan(&rec, []string{"A", "B", "C", "D", "hhhh", "1234"})
	if err != nil {
		t.Errorf("%+v", err)
	}

	fmt.Printf("%+v\n", rec)
}

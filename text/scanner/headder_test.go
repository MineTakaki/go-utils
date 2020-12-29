package scanner

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestWithHeadder(t *testing.T) {
	type testT1 struct {
		Col1 string         `headder:"Column1"`
		Col2 string         `headder:"Column2"`
		Col3 string         `json:"hoge" headder:"Column3"`
		Col4 string         `headder:"Column4"`
		Col5 sql.NullString `headder:"Column5"`
	}

	headers := []string{
		"Column1",
		"Column2",
		"Column3",
		"Column4",
		"Column5",
	}

	rec := testT1{}

	scan, err := WithHeadder(&rec, "headder", headers, nil)
	if err != nil {
		t.Errorf("%+v", err)
	}

	err = scan.Scan(&rec, []string{"A", "B", "C", "D", "hhhh"})
	if err != nil {
		t.Errorf("%+v", err)
	}

	t.Logf("%+v", rec)
	fmt.Printf("%+v\n", rec)

	return
}

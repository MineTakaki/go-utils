package decimal

import (
	"encoding/json"
	"testing"
)

func TestNullDecimlal_UnmarshallJson(t *testing.T) {
	var doc struct {
		Amount NullDecimal `json:"amount"`
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

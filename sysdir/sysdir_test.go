package sysdir

import "testing"

func TestDocuments(t *testing.T) {
	d, err := Documents()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	t.Logf("Documents = %s", d)
}

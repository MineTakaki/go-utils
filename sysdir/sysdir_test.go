package sysdir_test

import (
	"testing"

	"github.com/MineTakaki/go-utils/sysdir"
)

func TestDocuments(t *testing.T) {
	d, err := sysdir.Documents()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	if d == "" {
		t.Error("Documents() return is empty")
		return
	}
	t.Logf("Documents = %s", d)
}

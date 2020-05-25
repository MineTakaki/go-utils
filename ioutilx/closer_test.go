package ioutilx

import (
	"io"
	"testing"
)

func TestCloseHolder(t *testing.T) {
	var a, b, c io.Closer

	ch := NewCloserHolder(a, b, c)
	err := ch.Close()
	if err != nil {
		t.Errorf("%+v", err)
	}
}

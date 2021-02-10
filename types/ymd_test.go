package types

import (
	"testing"
	"time"
)

func TestYmdScan(t *testing.T) {

	var tm time.Time
	var ymd Ymd
	if err := ymd.Scan(tm); err != nil {
		t.Errorf("%+v", err)
	} else {
		t.Logf("%v", ymd)
	}

	tm = time.Now()
	if err := ymd.Scan(tm); err != nil {
		t.Errorf("%+v", err)
	} else {
		t.Logf("%v", ymd)
	}

}

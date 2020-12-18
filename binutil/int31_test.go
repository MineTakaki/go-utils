package binutil

import (
	"encoding/hex"
	"math"
	"testing"
)

func TestPutAndGetInt31(t *testing.T) {
	buf := make([]byte, 5)

	for i := 0; i <= math.MaxInt32; i++ {
		n := PutInt31(buf, i)
		d, l, err := GetInt31(buf[:n])
		if err != nil {
			t.Errorf("%d: %s\n%+v", i, hex.Dump(buf[:n]), err)
			break
		}
		if n != l {
			t.Errorf("length error %d, %d != %d, %s", i, n, l, hex.Dump(buf[:n]))
			break
		}
		if d != i {
			t.Errorf("value error %d != %d, %s", i, d, hex.Dump(buf[:n]))
			break
		}

	}
}

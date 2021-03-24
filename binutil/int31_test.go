package binutil

import (
	"encoding/hex"
	"errors"
	"io"
	"math"
	"sync"
	"testing"
)

func TestPutAndGetInt31(t *testing.T) {
	buf := make([]byte, 5)

	wg := sync.WaitGroup{}
	wg.Add(1)

	fnSkip := func(i int) bool {
		// 0/011 110/0 0111 0/011 110/0 1111 0/011 1110
		return (i & 0x3C73CF3E) != 0
	}

	pr, pw := io.Pipe()
	go func() {
		defer func() {
			if err := pr.Close(); err != nil {
				t.Errorf("%+v", err)
			}
			wg.Done()
		}()
		for i := uint(0); ; i++ {
			if fnSkip(int(i)) {
				continue
			}
			d, _, err := ReadInt31(pr)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					t.Errorf("%+v", err)
					return
				}
				break
			}
			if d != int(i) {
				t.Errorf("value error(ReadInt31) %d != %d", i, d)
			}
		}
	}()

	func() {
		defer func() {
			if err := pw.Close(); err != nil {
				t.Errorf("%+v", err)
			}
		}()
		cnt := 0
		for i := uint(0); i <= math.MaxInt32; i++ {
			if fnSkip(int(i)) {
				continue
			}
			cnt++
			n := PutInt31(buf, int(i))
			d, l, err := GetInt31(buf[:n])
			if err != nil {
				t.Errorf("%d: %s\n%+v", i, hex.Dump(buf[:n]), err)
				break
			}
			if n != l {
				t.Errorf("length error %d, %d != %d, %s", i, n, l, hex.Dump(buf[:n]))
				break
			}
			if d != int(i) {
				t.Errorf("value error %d != %d, %s", i, d, hex.Dump(buf[:n]))
				break
			}
			_, err = pw.Write(buf[:n])
			if err != nil {
				t.Errorf("%+v", err)
				return
			}
		}
		t.Logf("cnt: %d", cnt)
	}()

	wg.Wait()
}

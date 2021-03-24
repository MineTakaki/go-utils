package binutil

import (
	"encoding/hex"
	"errors"
	"io"
	"math"
	"sync"
	"testing"
)

func TestPutAndGetInt64(t *testing.T) {
	buf := make([]byte, 11)

	fnPow2 := func(n int) (d int64) {
		d = 1
		d <<= n
		return
	}

	data := func() (list []int64) {
		list = append(list, 0)
		list = append(list, math.MaxInt64)
		list = append(list, math.MinInt64)
		for i := 6; i < 63; i += 6 {
			d := fnPow2(i+1) - 1
			// 01/11 1111
			list = append(list, d)
			list = append(list, -d)

			// 01/01 1111
			x := fnPow2(i - 1)
			list = append(list, d-x)
			list = append(list, -(d - x))

			// 00/11 1111
			d >>= 1
			list = append(list, d)
			list = append(list, -d)

			// 00/01 1111
			d >>= 1
			list = append(list, d)
			list = append(list, -d)
		}
		return
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	pr, pw := io.Pipe()

	go func() {
		defer func() {
			if err := pr.Close(); err != nil {
				t.Errorf("%v", err)
			}
			wg.Done()
		}()
		for i := 0; ; i++ {
			d, _, err := ReadInt64(pr)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					t.Errorf("%d: %v", data[i], err)
					return
				}
				break
			}
			if data[i] != d {
				t.Errorf("value error %d != %d", data[i], d)
			}
		}
	}()

	func() {
		defer func() {
			if err := pw.Close(); err != nil {
				t.Errorf("%v", err)
			}
		}()
		cnt := 0
		for _, i := range data {
			cnt++
			n := PutInt64(buf, i)
			d, l, err := GetInt64(buf[:n])
			if err != nil {
				t.Errorf("%d: %s\n%+v", i, hex.Dump(buf[:n]), err)
				return
			}
			if n != l {
				t.Errorf("length error %d, %d != %d, %s", i, n, l, hex.Dump(buf[:n]))
				return
			}
			if d != i {
				t.Errorf("value error %d != %d, %s", i, d, hex.Dump(buf[:n]))
				return
			}
			_, err = pw.Write(buf[:n])
			if err != nil {
				t.Errorf("%v", err)
				return
			}
			//fmt.Printf("data: %064b, %d\n", uint64(i), i)
		}
		t.Logf("cnt: %d", cnt)
	}()

	wg.Wait()
}

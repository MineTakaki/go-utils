package binutil

import (
	"io"
	"math"

	"github.com/pkg/errors"
)

//PutInt64 int64をutf8 like なエンコードで[]byteに書き出します
// 最大11byteの領域が必要です
func PutInt64(b []byte, d int64) int {
	sign := false
	adjust := false
	if d < 0 {
		sign = true
		// マイナスの方が範囲が1広いので調整を行います
		if d == math.MinInt64 {
			adjust = true
			d++
		}
		d = -d
	}

	b[0] = byte(d & 0x3f)
	if sign {
		b[0] += 0x40
	}
	if d >>= 6; d == 0 {
		return 1
	}

	b[0] += 0x80
	for n := 1; ; n++ {
		b[n] = byte(d & 0x3f)
		if d >>= 6; d == 0 {
			b[n] += 0x80
			if adjust {
				b[n] += 0x8
			}
			return n + 1
		}
		b[n] += 0xC0
	}
}

//GetInt64 PutInt64()で書き出した[]byteからint64を復元します
func GetInt64(b []byte) (d int64, n int, err error) {
	if len(b) == 0 {
		err = errors.Errorf("decode error(no data)")
		return
	}
	n++
	sign := (b[0] & 0x40) != 0
	if (b[0] & 0x80) == 0 {
		d = int64(b[0] & 0x3f)
		if sign {
			d = -d
		}
		return
	}
	for ; n < 11; n++ {
		if len(b) <= n {
			err = errors.Errorf("decode error(not enouph bytes)")
			return
		}
		if (b[n] & 0x80) == 0 {
			err = errors.Errorf("decode error(data mark error)")
			return
		}
		if (b[n] & 0x40) == 0 {
			adjust := false
			x := n
			if n == 10 {
				d = int64(b[10] & 0x7)
				adjust = (b[10] & 0x8) == 0x8
				x--
			}
			for ; x >= 0; x-- {
				d = (d << 6) + int64(b[x]&0x3f)
			}
			if sign {
				d = -d
				if adjust {
					d--
				}
			}
			n++
			return
		}
	}
	err = errors.Errorf("decode error(end mark not found)")
	return
}

//ReadInt64 io.ReaderからPutInt64()で書き出したint64を読み取ります
func ReadInt64(r io.Reader) (d int64, n int, err error) {
	b := make([]byte, 11)
	if _, err = r.Read(b[0:1]); err != nil {
		err = errors.WithStack(err)
		return
	}
	n++
	sign := (b[0] & 0x40) != 0
	if (b[0] & 0x80) == 0 {
		d = int64(b[0] & 0x3f)
		if sign {
			d = -d
		}
		return
	}
	for ; n < 11; n++ {
		if _, err = r.Read(b[n : n+1]); err != nil {
			err = errors.WithStack(err)
			return
		}
		if (b[n] & 0x80) == 0 {
			err = errors.Errorf("decode error(data mark error)")
			return
		}
		if (b[n] & 0x40) == 0 {
			adjust := false
			x := n
			if n == 10 {
				d = int64(b[10] & 0x7)
				adjust = (b[10] & 0x8) == 0x8
				x--
			}
			for ; x >= 0; x-- {
				d = (d << 6) + int64(b[x]&0x3f)
			}
			if sign {
				d = -d
				if adjust {
					d--
				}
			}
			n++
			return
		}
	}
	err = errors.Errorf("decode error(end mark not found)")
	return
}

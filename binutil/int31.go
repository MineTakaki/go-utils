package binutil

import (
	"io"

	"github.com/pkg/errors"
)

//PutInt31 intをutf8 like なエンコードで[]byteに書き出します
// マイナス値は対応しません
func PutInt31(b []byte, d int) int {
	if d < 0 {
		d = 0
	}
	b[0] = byte(d & 0x7f)
	if d >>= 7; d == 0 {
		return 1
	}

	b[0] += byte(0x80)
	c := 1
	for func() bool {
		b[c] = byte((d & 0x3f) + 0x80)
		d >>= 6
		return d != 0
	}() {
		b[c] += 0x40
		c++
	}
	return c + 1
}

//GetInt31 PutInt31()で書き出した[]byteからintを復元します
func GetInt31(b []byte) (int, int, error) {
	switch {
	case len(b) == 0:
	case (b[0] & 0x80) == 0:
		return int(b[0] & 0x7f), 1, nil //7
	case len(b) < 2:
	case (b[1] & 0x80) == 0:
	case (b[1] & 0x40) == 0:
		return int(b[0]&0x7f) + int(b[1]&0x3f)<<7, 2, nil //13
	case len(b) < 3:
	case (b[2] & 0x80) == 0:
	case (b[2] & 0x40) == 0:
		return int(b[0]&0x7f) + (int(b[1]&0x3f)+(int(b[2]&0x3f)<<6))<<7, 3, nil //19
	case len(b) < 4:
	case (b[3] & 0x80) == 0:
	case (b[3] & 0x40) == 0:
		return int(b[0]&0x7f) + (int(b[1]&0x3f)+(int(b[2]&0x3f)+int(b[3]&0x3f)<<6)<<6)<<7, 4, nil //25
	case len(b) < 5:
	case (b[4] & 0x80) == 0:
	case (b[4] & 0x40) == 0:
		return int(b[0]&0x7f) + (int(b[1]&0x3f)+(int(b[2]&0x3f)+(int(b[3]&0x3f)+int(b[4]&0x3f)<<6)<<6)<<6)<<7, 5, nil //31
	}
	return 0, 0, errors.Errorf("decode error")
}

func ReadInt31(r io.Reader) (d, n int, err error) {
	b := make([]byte, 1)
	if _, err = r.Read(b); err != nil {
		err = errors.WithStack(err)
		return
	}
	n++
	d = int(b[0] & 0x7f)
	if (b[0] & 0x80) == 0 {
		return
	}

	if _, err = r.Read(b); err != nil {
		err = errors.WithStack(err)
		return
	}
	n++
	d += int(b[0]&0x3f) << 7
	if (b[0] & 0x80) == 0 {
		err = errors.Errorf("decode error")
		return
	}
	if (b[0] & 0x40) == 0 {
		return
	}

	if _, err = r.Read(b); err != nil {
		err = errors.WithStack(err)
		return
	}
	n++
	d += int(b[0]&0x3f) << 13
	if (b[0] & 0x80) == 0 {
		err = errors.Errorf("decode error")
		return
	}
	if (b[0] & 0x40) == 0 {
		return
	}

	if _, err = r.Read(b); err != nil {
		err = errors.WithStack(err)
		return
	}
	n++
	d += int(b[0]&0x3f) << 19
	if (b[0] & 0x80) == 0 {
		err = errors.Errorf("decode error")
		return
	}
	if (b[0] & 0x40) == 0 {
		return
	}

	if _, err = r.Read(b); err != nil {
		err = errors.WithStack(err)
		return
	}
	n++
	d += int(b[0]&0x3f) << 25
	if (b[0] & 0x80) == 0 {
		err = errors.Errorf("decode error")
		return
	}
	if (b[0] & 0x40) == 0 {
		return
	}

	err = errors.Errorf("decode error")
	return
}

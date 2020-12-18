package binutil

import "github.com/pkg/errors"

//PutInt31 intをutf8 like なエンコードで[]byteに書き出します
// マイナス値は対応しません
func PutInt31(b []byte, d int) int {
	u := uint32(d)
	b[0] = byte(u & 0x7f)
	if u >>= 7; u == 0 {
		return 1
	}

	b[0] += byte(0x80)
	c := 1
	for func() bool {
		b[c] = byte((u & 0x3f) + 0x80)
		u >>= 6
		return u != 0
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

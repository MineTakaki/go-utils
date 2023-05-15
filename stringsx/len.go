package stringsx

import "unicode/utf8"

// LenW 全角を2、半角を1として string の長さを取得します
func LenW(s string) (x int) {
	ToBytes(s, func(b []byte) {
		for s := 0; s < len(b); {
			_, n := utf8.DecodeRune(b[s:])
			if n <= 1 {
				s++
				x++
				continue
			}
			s += n
			x += 2
		}
	})
	return
}

package stringsx

import (
	"sort"
	"strings"
)

func UniqueStrings(n int, fn func(i int) string) []string {
	var buf []string
	for i := 0; i < n; i++ {
		s := fn(i)
		x, found := sort.Find(len(buf), func(i int) int { return strings.Compare(buf[i], s) })
		if found {
			continue
		}
		buf = append(buf, "")
		if len(buf) < n {
			copy(buf[x+1:], buf[x:])
		}
		buf[x] = s
	}
	return buf
}

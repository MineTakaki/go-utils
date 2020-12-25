package stringsx

import (
	"bufio"
	"strconv"
	"strings"
)

//Ln 行番号を付けて返します
func Ln(s string) string {
	sc := bufio.NewScanner(strings.NewReader(s))
	sc.Split(bufio.ScanLines)
	sb := strings.Builder{}
	n := 0
	for sc.Scan() {
		n++
		sb.WriteString(strconv.Itoa(n) + ":" + sc.Text() + "\n")
	}
	return sb.String()
}

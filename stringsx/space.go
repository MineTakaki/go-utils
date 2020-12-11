package stringsx

import (
	"strings"
	"unicode"
)

//ReduceSpace 複数の空白をひとまとめにします
func ReduceSpace(str string) string {
	//重複する空白を削除、ついでに前後の空白もトリムします
	var sb strings.Builder
	var flg bool
	for _, r := range str {
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			flg = true
			continue
		}
		if flg {
			flg = false
			if sb.Len() != 0 {
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune(r)
	}
	return sb.String()
}

package stringsx

import "github.com/MineTakaki/go-utils/runesx"

//Reverse 文字列を反転します
func Reverse(str string) string {
	return string(runesx.Reverse([]rune(str)))
}

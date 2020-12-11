package runesx

func minN(a, b int) int {
	if a > b {
		return b
	}
	return a
}

//MatchLengh 先頭から一致する長さを取得します
func MatchLengh(a, b []rune) (n int) {
	m := minN(len(a), len(b))
	for i := 0; i < m && a[i] == b[i]; i++ {
		n++
	}
	return
}

//MatchLenghLast 末尾から一致する長さを取得します
func MatchLenghLast(a, b []rune) (n int) {
	for i, j := len(a)-1, len(b)-1; i >= 0 && j >= 0 && a[i] == b[j]; {
		n++
		i--
		j--
	}
	return
}

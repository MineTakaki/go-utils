package merge

//MergeU 大小比較によるマージ処理のヘルパー(キー値がユニークな場合のみ限定)
func MergeU(
	cntA, cntB int,
	cmp func(i, j int) int,
	merge func(cr, i, j int) error,
) (err error) {
	i, j := 0, 0
	for {
		if i >= cntA {
			if j >= cntB {
				break
			}
			for ; j < cntB; j++ {
				if err = merge(1, i, j); err != nil {
					return
				}
			}
			break
		}
		if j >= cntB {
			for ; i < cntA; i++ {
				if err = merge(-1, i, j); err != nil {
					return
				}
			}
			break
		}
		cr := cmp(i, j)
		if err = merge(cr, i, j); err != nil {
			return
		}
		switch {
		default:
			i++
			j++
		case cr < 0:
			i++
		case cr > 0:
			j++
		}
	}
	return
}

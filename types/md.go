package types

import (
	"fmt"

	"github.com/MineTakaki/go-utils/conv"
)

type (
	//Md MMdd形式で月日を表す整数型
	Md int
)

//ToMd 月日からMd型に変換します
func ToMd(m, d int) (md Md, err error) {
	_, err = ValidateMd(m, d)
	if err != nil {
		return
	}
	md = Md(m*100 + d)
	return
}

//ParseMd Md型に変換します
func ParseMd(i interface{}) (md Md, err error) {
	err = md.Scan(i)
	if err == nil {
		_, err = md.Validate()
	}
	return
}

//ParseMd2 Md型に変換します
func ParseMd2(i interface{}, err *error) (md Md) {
	var e error
	md, e = ParseMd(i)
	if err != nil {
		*err = e
	}
	return
}

//String string型変換
func (md Md) String() string {
	if md == 0 {
		return ""
	}
	return fmt.Sprintf("%04d", md)
}

//Validate 年月が正しいか確認します
func (md Md) Validate() (bool, error) {
	if md == 0 {
		return true, nil
	}

	m, d := md.Part()
	return ValidateMd(m, d)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

//Month 月を取得します
func (md Md) Month() int {
	return int(md) / 100
}

//Day 日を取得します
func (md Md) Day() int {
	return abs(int(md)) % 100
}

//Part 月日の要素を取得します
func (md Md) Part() (m, d int) {
	m = int(md) / 100
	d = abs(int(md)) % 100
	return
}

//Prev 前日の値を取得します（うるう年判定は行いません）
func (md Md) Prev() Md {
	return md.Add(0, -1)
}

//Next 翌日の値を取得します（うるう年判定は行いません）
func (md Md) Next() Md {
	return md.Add(0, 1)
}

var _days = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

//AdjustMonth 月の加減算後の正しい月を取得します
func AdjustMonth(m int) (int, int) {
	if m < 1 {
		dy := m/12 - 1
		return m - dy*12, dy
	} else if m > 12 {
		dy := (m - 1) / 12
		return m - dy*12, dy
	}
	return m, 0
}

//AdjustDay 日の加減算後の正しい日を取得します
func AdjustDay(m, d int) (int, int) {
	if m < 1 || m > 12 {
		m, _ = AdjustMonth(m)
	}

	dm := 0
	if d < 1 {
		for d < 1 {
			dm--
			if m--; m < 1 {
				m = 12
			}
			maxd := _days[m-1]
			d = d + maxd
		}
		return d, dm
	}

	for {
		maxd := _days[m-1]
		if d <= maxd {
			break
		}
		d = d - maxd
		if m++; m > 12 {
			m = 1
		}
		dm++
	}
	return d, dm
}

//Add 月、日を加算します（減算はマイナス値を引数にセットします）
func (md Md) Add(dm, dd int) Md {
	m, d := md.Part()

	//月の状態を正常化します
	m, _ = AdjustMonth(m)

	//日を計算します
	var x int
	d, x = AdjustDay(m, d+dd)

	//月を計算します
	m, _ = AdjustMonth(m + dm + x)

	return Md(m*100 + d)
}

//Scan 文字列から月日を読み取ります
func (md *Md) Scan(i interface{}) (err error) {
	if conv.IsEmpty(i) {
		*md = 0
		return nil
	}
	if n, ok := conv.Int(i); ok {
		*md = Md(n)
		return
	}
	return ErrValidate
}

//Between 二つの日付の間に入るか判定します
func (md Md) Between(f, t Md) bool {
	if md == 0 || f == 0 || t == 0 {
		return false
	}
	if f > t {
		return md >= f || md <= t
	}
	return f <= md && md >= t
}

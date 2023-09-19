package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/MineTakaki/go-utils/conv"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

type (
	//Md MMdd形式で月日を表す整数型
	Md int

	//MdSlice Md型のスライス
	MdSlice []Md
)

func (md MdSlice) Len() int           { return len(md) }
func (md MdSlice) Less(i, j int) bool { return md[i] < md[j] }
func (md MdSlice) Swap(i, j int)      { md[i], md[j] = md[j], md[i] }

// ToMd 月日からMd型に変換します
func ToMd(m, d int) (md Md, err error) {
	_, err = ValidateMd(m, d)
	if err != nil {
		return
	}
	md = Md(m*100 + d)
	return
}

// ParseMd Md型に変換します
func ParseMd(i interface{}) (md Md, err error) {
	err = md.Scan(i)
	if err == nil {
		_, err = md.Validate()
	}
	return
}

// ParseMd2 Md型に変換します
func ParseMd2(i interface{}, err *error) (md Md) {
	var e error
	md, e = ParseMd(i)
	if err != nil {
		*err = e
	}
	return
}

// String string型変換
func (md Md) String() string {
	if md == 0 {
		return ""
	}
	return fmt.Sprintf("%04d", md)
}

// FormatMd MD形式でstring型に整形して変換します
func (md Md) FormatMd(sep string, zeroSuppress bool) string {
	if md == 0 {
		return ""
	}
	m, d := md.Part()
	sb := strings.Builder{}
	sb.Grow(len(sep) + 4)
	if zeroSuppress {
		sb.WriteString(strconv.Itoa(m))
		sb.WriteString(sep)
		sb.WriteString(strconv.Itoa(d))
		return sb.String()
	}
	sb.WriteString(fillZero2(m))
	sb.WriteString(sep)
	sb.WriteString(fillZero2(d))
	return sb.String()
}

// Value driver.Valuerインターフェイスの実装
func (md Md) Value() (driver.Value, error) {
	if md == 0 {
		return nil, nil
	}
	return int64(md), nil
}

// UnmarshalJSON json.Unmarshalerインターフェイスの実装
func (md *Md) UnmarshalJSON(b []byte) (err error) {
	var s interface{}
	if err = json.Unmarshal(b, &s); err != nil {
		err = errors.WithStack(err)
		return
	}
	var x Md
	if x, err = ParseMd(s); err != nil {
		return
	}
	*md = x
	return
}

// MarshalJSON json.Marshalerの実装
func (md *Md) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(*md))
}

// MarshalLogObject zapcore.ObjectMarshalerの実装
func (md *Md) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("md", md.String())
	return nil
}

// Validate 年月が正しいか確認します
func (md Md) Validate() (bool, error) {
	if md == 0 {
		return true, nil
	}

	m, d := md.Part()
	return ValidateMd(m, d)
}

// Month 月を取得します
func (md Md) Month() int {
	return int(md) / 100
}

// Day 日を取得します
func (md Md) Day() int {
	return abs(int(md)) % 100
}

// Part 月日の要素を取得します
func (md Md) Part() (m, d int) {
	m = int(md) / 100
	d = abs(int(md)) % 100
	return
}

// Parts 月日の要素を配列で取得します
func (md Md) Parts() []int {
	v := make([]int, 2)
	v[0], v[1] = md.Part()
	return v
}

// Prev 前日の値を取得します（うるう年判定は行いません）
func (md Md) Prev() Md {
	return md.Add(0, -1)
}

// Next 翌日の値を取得します（うるう年判定は行いません）
func (md Md) Next() Md {
	return md.Add(0, 1)
}

var _days = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

// AdjustMonth 月の加減算後の正しい月を取得します
func AdjustMonth(y, m int) (qy, qm int) {
	if m < 1 {
		dy := m/12 - 1
		qy = y + dy
		qm = m - dy*12
		return
	} else if m > 12 {
		dy := (m - 1) / 12
		qy = y + dy
		qm = m - dy*12
		return
	}
	qy, qm = y, m
	return
}

// AdjustDay 日の加減算後の正しい日を取得します
func AdjustDay(y, m, d int) (int, int, int) {
	if m < 1 || m > 12 {
		y, m = AdjustMonth(y, m)
	}

	if d < 1 {
		for d < 1 {
			if m--; m < 1 {
				m = 12
				y--
			}
			d += LastDay(y, m)
		}
		return y, m, d
	}

	for {
		maxd := LastDay(y, m)
		if d <= maxd {
			break
		}
		d -= maxd
		if m++; m > 12 {
			m = 1
			y++
		}
	}
	return y, m, d
}

// AdjustMd 月日の加減算後の正しい日を取得します(うるう年は考慮しない)
func AdjustMd(m, d int) (int, int) {
	if m < 1 {
		m -= (m/12 - 1) * 12
	} else if m > 12 {
		m -= ((m - 1) / 12) * 12
	}

	if d < 1 {
		for d < 1 {
			if m--; m < 1 {
				m = 12
			}
			d += _days[m-1]
		}
		return m, d
	}

	for {
		maxd := _days[m-1]
		if d <= maxd {
			break
		}
		d -= maxd
		if m++; m > 12 {
			m = 1
		}
	}
	return m, d
}

// Add 月、日を加算します（減算はマイナス値を引数にセットします）
func (md Md) Add(dm, dd int) Md {
	if md == 0 {
		return 0
	}
	m, d := md.Part()

	//状態を正常化します
	m, d = AdjustMd(m, d+dd)

	return Md(m*100 + d)
}

// Adjust 月日を正しい形式に訂正します
func (md Md) Adjust() Md {
	return md.Add(0, 0)
}

// Scan 文字列から月日を読み取ります
func (md *Md) Scan(i interface{}) (err error) {
	if conv.IsEmpty(i) {
		*md = 0
		return nil
	}
	if n, ok := conv.Int(i); ok {
		*md = Md(n)
		return
	}
	return errors.WithStack(ErrValidate)
}

// Between 二つの日付の間に入るか判定します
func (md Md) Between(f, t Md) bool {
	if md == 0 || f == 0 || t == 0 {
		return false
	}
	if f > t {
		return md >= f || md <= t
	}
	return f <= md && md >= t
}

// Compare Md同志を比較します
func (md Md) Compare(o Md) int {
	if md < o {
		return -1
	}
	if md > o {
		return 1
	}
	return 0
}

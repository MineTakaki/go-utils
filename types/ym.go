package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MineTakaki/go-utils/conv"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

type (
	//Ym yyyyMM形式で年月を表す整数型
	Ym int

	//YmSlice Ym型のスライス
	YmSlice []Ym
)

// ErrValidate 値が適切でない
var ErrValidate = errors.New("validate error")

// ErrUnkownType 知らない型が指定されました
var ErrUnkownType = errors.New("unkown type")

func (ym YmSlice) Len() int           { return len(ym) }
func (ym YmSlice) Less(i, j int) bool { return ym[i] < ym[j] }
func (ym YmSlice) Swap(i, j int)      { ym[i], ym[j] = ym[j], ym[i] }

// ToYm 年月からYm型に変換します
func ToYm(y, m int) (ym Ym, err error) {
	_, err = ValidateYm(y, m)
	if err != nil {
		return
	}
	ym = Ym(y*100 + m)
	return
}

// ParseYm Ym型に変換します
func ParseYm(i interface{}) (ym Ym, err error) {
	err = ym.Scan(i)
	if err == nil {
		_, err = ym.Validate()
	}
	return
}

// ParseYm2 Ym型に変換します
func ParseYm2(i interface{}, err *error) (ym Ym) {
	var e error
	ym, e = ParseYm(i)
	if err != nil {
		*err = e
	}
	return
}

// String string型変換
func (ym Ym) String() string {
	if ym == 0 {
		return ""
	}
	return fmt.Sprintf("%06d", ym)
}

// FormatYm YM形式でstring型に整形して変換します
func (ym Ym) FormatYm(sep string, zeroSuppress bool) string {
	if ym == 0 {
		return ""
	}
	y, m := ym.Part()
	sb := strings.Builder{}
	sb.Grow(len(sep) + 6)
	if zeroSuppress {
		sb.WriteString(strconv.Itoa(y))
		sb.WriteString(sep)
		sb.WriteString(strconv.Itoa(m))
		return sb.String()
	}
	sb.WriteString(fillZero4(y))
	sb.WriteString(sep)
	sb.WriteString(fillZero2(m))
	return sb.String()
}

// Validate 年月が正しいか確認します
func (ym Ym) Validate() (bool, error) {
	if ym == 0 {
		return true, nil
	}

	y, m := ym.Part()
	return ValidateYm(y, m)
}

// Year 年を取得します
func (ym Ym) Year() int {
	return int(ym) / 100
}

// Month 月を取得します
func (ym Ym) Month() int {
	return int(ym) % 100
}

// Part 年月の要素を取得します
func (ym Ym) Part() (y, m int) {
	y = int(ym) / 100
	m = int(ym) % 100
	return
}

// Parts 年月の要素を配列で取得します
func (ym Ym) Parts() []int {
	v := make([]int, 2)
	v[0], v[1] = ym.Part()
	return v
}

// Prev 一か月前の値を取得します
func (ym Ym) Prev() Ym {
	return ym.Add(0, -1)
}

// Next 一か月後の値を取得します
func (ym Ym) Next() Ym {
	return ym.Add(0, 1)
}

// Add 年、月を加算します（減算はマイナス値を引数にセットします）
func (ym Ym) Add(dy, dm int) Ym {
	if ym == 0 {
		return 0
	}
	y, m := ym.Part()
	y, m = AdjustMonth(y+dy, m+dm)
	return Ym(y*100 + m)
}

// Ymd 日を指定してYmd型を取得します
func (ym Ym) Ymd(d int) Ymd {
	if ym == 0 {
		return 0
	}
	xy, xm := ym.Part()
	if d < 1 {
		d = 1
	} else if lday := LastDay(xy, xm); d > lday {
		d = lday
	}
	return Ymd(xy*10000 + xm*100 + d)
}

// GoTime go言語のTime型を取得します
func (ym Ym) GoTime() (tm time.Time) {
	if ym != 0 {
		tm = time.Date(ym.Year(), time.Month(ym.Month()), 1, 0, 0, 0, 0, time.Local)
	}
	return
}

// Term From～To
func (ym Ym) Term() (fm, to Ymd) {
	if ym == 0 {
		return
	}
	fm = Ymd(ym*100 + 1)
	y, m := ym.Part()
	to = Ymd(int(ym)*100 + LastDay(y, m))
	return
}

// First 月初の年月日を取得します
func (ym Ym) First() Ymd {
	return Ymd(ym*100 + 1)
}

// Last 月末の年月日を取得します
func (ym Ym) Last() Ymd {
	y, m := ym.Part()
	return Ymd(int(ym)*100 + LastDay(y, m))
}

// BetweenMonth 第1引数の月が第2,3引数の月の間に該当するか判定します
func BetweenMonth(m, m1, m2 int) bool {
	if m == 0 {
		return false
	}
	if m1 <= m2 {
		return m >= m1 && m <= m2
	}
	return m >= m1 || m <= m2
}

// BetweenMonth 2つの月の間に該当するか判定します。m1 > m2の場合は年を跨ぐ範囲として扱います
func (ym Ym) BetweenMonth(m1, m2 int) bool {
	if ym == 0 {
		return false
	}
	return BetweenMonth(ym.Month(), m1, m2)
}

// Scan 年月を読み取ります
func (ym *Ym) Scan(i interface{}) (err error) {
	if conv.IsEmpty(i) {
		*ym = 0
		return nil
	}
	if tm, ok := i.(time.Time); ok {
		*ym = Ym(tm.Year()*100 + int(tm.Month()))
		return nil
	}
	if tm, ok := i.(*time.Time); ok {
		*ym = Ym(tm.Year()*100 + int(tm.Month()))
		return nil
	}
	if n, ok := conv.Int(i); ok {
		*ym = Ym(n)
		return
	}
	if s, ok := i.(string); ok {
		for _, layout := range []string{"2006-01", "2006/01"} {
			if tm, err := time.Parse(layout, s); err == nil {
				*ym = Ym(tm.Year()*100 + int(tm.Month()))
				return nil
			}
		}
	}
	return errors.WithStack(ErrValidate)
}

// Value driver.Valuerインターフェイスの実装
func (ym Ym) Value() (driver.Value, error) {
	if ym == 0 {
		return nil, nil
	}
	return int64(ym), nil
}

// UnmarshalJSON json.Unmarshalerインターフェイスの実装
func (ym *Ym) UnmarshalJSON(b []byte) (err error) {
	var s interface{}
	if err = json.Unmarshal(b, &s); err != nil {
		err = errors.WithStack(err)
		return
	}
	var x Ym
	if x, err = ParseYm(s); err != nil {
		return
	}
	*ym = x
	return
}

// MarshalJSON json.Marshalerの実装
func (ym *Ym) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(*ym))
}

// MarshalLogObject zapcore.ObjectMarshalerの実装
func (ym *Ym) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("ym", ym.String())
	return nil
}

// CsvFormat CSV出力用のstring型変換
func (ym Ym) CsvFormat() string {
	if ym == 0 {
		return ""
	}
	return strconv.Itoa(int(ym))
}

// IsLeapYear うるう年判定
func IsLeapYear(y int) bool {
	return ((y%4) == 0 && (y%100) != 0) || (y%400) == 0
}

// Min 指定した年月と比較して小さい値を返します
func (ym Ym) Min(o Ym) Ym {
	if ym == 0 {
		return o
	}
	if o != 0 && ym > o {
		return o
	}
	return ym
}

// Max 指定した年月と比較して大きい値を返します
func (ym Ym) Max(o Ym) Ym {
	if ym == 0 {
		return o
	}
	if o != 0 && ym < o {
		return o
	}
	return ym
}

// Compare Ym同志を比較します
func (ym Ym) Compare(o Ym) int {
	if ym < o {
		return -1
	}
	if ym > o {
		return 1
	}
	return 0
}

// TermYear 区切りを指定して年度を取得します
func (ym Ym) TermYear(start int) int {
	if ym == 0 {
		return 0
	}
	y := ym.Year()
	if ym.Month() < start {
		y--
	}
	return y
}

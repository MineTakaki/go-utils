package types

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"time"

	"github.com/MineTakaki/go-utils/conv"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

type (
	// Ymdhms yyyyMMdd形式で年月日時分秒を表す整数型
	Ymdhms int64

	//YmdhmsSlice Ymdhms型のスライス
	YmdhmsSlice []Ymdhms
)

func (yh YmdhmsSlice) Len() int           { return len(yh) }
func (yh YmdhmsSlice) Less(i, j int) bool { return yh[i] < yh[j] }
func (yh YmdhmsSlice) Swap(i, j int)      { yh[i], yh[j] = yh[j], yh[i] }

// ToYmdhms 年月日時分秒からYmdhms型に変換します
func ToYmdhms(y, m, d, h, n, s int) (yh Ymdhms, err error) {
	_, err = ValidateYmdhms(y, m, d, h, n, s)
	if err == nil {
		yh = Ymdhms(((((int64(y)*100+int64(m))*100+int64(d))*100+int64(h))*100+int64(n))*100 + int64(s))
	}
	return
}

// ParseYmdhms Ymdhms型に変換します
func ParseYmdhms(i interface{}) (yh Ymdhms, err error) {
	if err = yh.Scan(i); err != nil {
		return
	}
	_, err = yh.Validate()
	return
}

// ParseYmdhms2 Ymdhms型に変換します
func ParseYmdhms2(i interface{}, err *error) (yh Ymdhms) {
	var e error
	yh, e = ParseYmdhms(i)
	if err != nil {
		*err = e
	}
	return
}

// YmdhmsFromGoTime time.TimeからYmdhmsに変換します
func YmdhmsFromGoTime(tm time.Time) Ymdhms {
	return Ymdhms(((((int64(tm.Year())*100+int64(tm.Month()))*100+int64(tm.Day()))*100+int64(tm.Hour()))*100+int64(tm.Minute()))*100 + int64(tm.Second()))
}

// YmdhmsNow 現在の日付（ローカル）を取得します
func YmdhmsNow() Ymdhms {
	return YmdhmsFromGoTime(time.Now())
}

// String string型変換
func (yh Ymdhms) String() string {
	if yh == 0 {
		return ""
	}
	return ItoaZeroFilled64(int64(yh), 14)
}

// CsvFormat CSV出力用のstring型変換
func (yh Ymdhms) CsvFormat() string {
	if yh == 0 {
		return ""
	}
	return strconv.FormatInt(int64(yh), 10)
}

// FormatYmdhms string型に整形して変換します
func (yh Ymdhms) FormatYmdhms(sep, sep2 string, zeroSuppress bool) string {
	if yh == 0 {
		return ""
	}
	return yh.Ymd().FormatYmd(sep, zeroSuppress) + " " + yh.Hms().FormatHms(sep2, zeroSuppress)
}

// Scan 年月日時分秒を読み取ります
func (yh *Ymdhms) Scan(i interface{}) error {
	if conv.IsEmpty(i) {
		*yh = 0
		return nil
	}

	if tm, ok := i.(time.Time); ok {
		*yh = YmdhmsFromGoTime(tm)
		return nil
	}
	if tm, ok := i.(*time.Time); ok {
		*yh = YmdhmsFromGoTime(*tm)
		return nil
	}
	if n, ok := conv.Int64(i); ok {
		*yh = Ymdhms(n)
		return nil
	}
	if s, ok := i.(string); ok {
		for _, layout := range []string{"2006-01-02 15:04:05", "2006/01/02 15:04:05", "2006-1-2 15:04:05", "2006/1/2 15:04:05"} {
			if tm, err := time.Parse(layout, s); err == nil {
				*yh = YmdhmsFromGoTime(tm)
				return nil
			}
		}
	}
	return errors.WithStack(ErrValidate)
}

// Value driver.Valuerインターフェイスの実装
func (yh Ymdhms) Value() (driver.Value, error) {
	if yh == 0 {
		return nil, nil
	}
	return int64(yh), nil
}

// UnmarshalJSON json.Unmarshalerインターフェイスの実装
func (yh *Ymdhms) UnmarshalJSON(b []byte) (err error) {
	var s interface{}
	if err = json.Unmarshal(b, &s); err != nil {
		err = errors.WithStack(err)
		return
	}
	var x Ymdhms
	if x, err = ParseYmdhms(s); err != nil {
		return
	}
	*yh = x
	return
}

// MarshalJSON json.Marshalerの実装
func (yh *Ymdhms) MarshalJSON() ([]byte, error) {
	return json.Marshal(yh.String())
}

// MarshalLogObject zapcore.ObjectMarshalerの実装
func (yh *Ymdhms) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("ymdhms", yh.String())
	return nil
}

// Validate 年月日時分秒が正しいか確認します
func (yh Ymdhms) Validate() (bool, error) {
	if yh == 0 {
		return true, nil
	}
	y, m, d, h, n, s := yh.Part()
	return ValidateYmdhms(y, m, d, h, n, s)
}

// Part 年月日時分秒の要素を取得します
func (yh Ymdhms) Part() (y, m, d, h, n, s int) {
	d = int(yh) % 100
	x := int(yh) / 100
	m = x % 100
	y = x / 100
	return
}

// Parts 年月日時分秒の要素を配列で取得します
func (yh Ymdhms) Parts() []int {
	v := make([]int, 6)
	v[0], v[1], v[2], v[3], v[4], v[5] = yh.Part()
	return v
}

// GoTime go言語のTime型に変換します
func (yh Ymdhms) GoTime() (tm time.Time) {
	if yh != 0 {
		y, m, d, h, n, s := yh.Part()
		h, n, s = AdjustHms(h, n, s)
		if h > 23 {
			d += h / 24
			h %= 24
		}
		y, m, d = AdjustDay(y, m, d)
		tm = time.Date(y, time.Month(m), d, h, n, s, 0, time.Local)
	}
	return
}

// Add 年、月、日、時、分、秒を加算します（減算はマイナス値を引数にセットします）
func (yh Ymdhms) Add(y, m, d, h, n, s int) Ymdhms {
	if yh == 0 {
		return 0
	}

	h, n, s = yh.Hms().Add(h, n, s).Part()
	if h < 0 {
		h = abs(h)
		d -= (h / 24) + 1
		h = 24 - h%24
	} else if h > 23 {
		d += h / 24
		h %= 24
	}

	y, m, d = yh.Ymd().Add(y, m, d).Part()

	return Ymdhms(((((int64(y)*100+int64(m))*100+int64(d))*100+int64(h))*100+int64(n))*100 + int64(s))
}

// SetYmd Ymd型の値をセットします
//
//	負数の場合や日付時刻としての整合性は考慮しません
func (yh Ymdhms) SetYmd(ymd Ymd) Ymdhms {
	return Ymdhms(int64(ymd)*1000000 + abs64(int64(yh))%1000000)
}

// SetHms Hms型の値をセットします
//
//	負数の場合や日付時刻としての整合性は考慮しません
func (yh Ymdhms) SetHms(hms Hms) Ymdhms {
	return Ymdhms(int64(yh.Ymd())*1000000 + abs64(int64(hms))%1000000)
}

// SetDay 日を指定した値で置き換えます
func (yh Ymdhms) SetDay(d int) Ymdhms {
	if yh == 0 {
		return 0
	}
	return yh.SetYmd(yh.Ymd().SetDay(d))
}

// Prev 前日の値を取得します
func (yh Ymdhms) Prev() Ymdhms {
	return yh.Add(0, 0, 0, 0, 0, -1)
}

// Next 一日後の値を取得します
func (yh Ymdhms) Next() Ymdhms {
	return yh.Add(0, 0, 0, 0, 0, 1)
}

// Ymd 年月日をYmd型の値で取得します
func (yh Ymdhms) Ymd() Ymd {
	if yh == 0 {
		return 0
	}
	return Ymd(int(yh / 1000000))
}

// Hms 時分秒をHms型の値で取得します
func (yh Ymdhms) Hms() Hms {
	if yh == 0 {
		return 0
	}
	return Hms(int64(yh) % 1000000)
}

// Year 年を取得します
func (yh Ymdhms) Year() int {
	if yh == 0 {
		return 0
	}
	return yh.Ymd().Year()
}

// Month 月を取得します
func (yh Ymdhms) Month() int {
	if yh == 0 {
		return 0
	}
	return yh.Ymd().Month()
}

// Day 日を取得します
func (yh Ymdhms) Day() int {
	if yh == 0 {
		return 0
	}
	return yh.Ymd().Day()
}

// YearMonth 年月を取得します
//
// Deprecated
func (yh Ymdhms) YearMonth() Ym {
	if yh == 0 {
		return 0
	}
	return yh.Ymd().YearMonth()
}

// YearMonth 年月を取得します
func (yh Ymdhms) Ym() Ym {
	if yh == 0 {
		return 0
	}
	return yh.Ymd().YearMonth()
}

// MonthDay 月日を取得します
func (yh Ymdhms) MonthDay() Md {
	if yh == 0 {
		return 0
	}
	return yh.Ymd().MonthDay()
}

// MonthDay 月日を取得します
//
// deprecated
func (yh Ymdhms) Md() Md {
	if yh == 0 {
		return 0
	}
	return yh.Ymd().MonthDay()
}

// Between 二つの日付の間に入るか判定します
func (yh Ymdhms) Between(f, t Ymdhms) bool {
	if yh == 0 || f == 0 || t == 0 {
		return false
	}
	return f <= yh && yh <= t
}

// BetweenMonth 2つの月の間に該当するか判定します。m1 > m2の場合は年を跨ぐ範囲として扱います
func (yh Ymdhms) BetweenMonth(m1, m2 int) bool {
	if yh == 0 {
		return false
	}
	return BetweenMonth(yh.Month(), m1, m2)
}

// BetweenMonthDay 2つの月日の間に該当するか判定します。md1 > md2の場合は年を跨ぐ範囲として扱います
func (yh Ymdhms) BetweenMonthDay(md1, md2 Md) bool {
	return yh.MonthDay().Between(md1, md2)
}

// Min 指定した日付と比較して小さい値を返します
func (yh Ymdhms) Min(o Ymdhms) Ymdhms {
	if yh == 0 {
		return o
	}
	if o != 0 && yh > o {
		return o
	}
	return yh
}

// Max 指定した日付と比較して大きい値を返します
func (yh Ymdhms) Max(o Ymdhms) Ymdhms {
	if yh == 0 {
		return o
	}
	if o != 0 && yh < o {
		return o
	}
	return yh
}

// Days グレゴリウス暦1年1月1日からの経過日数を取得します
func (yh Ymdhms) Days() int {
	if yh == 0 {
		return 0
	}
	return yh.Ymd().Days()
}

// Compare Ymdhms同志を比較します
func (yh Ymdhms) Compare(o Ymdhms) int {
	if yh < o {
		return -1
	}
	if yh > o {
		return 1
	}
	return 0
}

// TermYear 区切りを指定して年度を取得します
func (yh Ymdhms) TermYear(start Md) int {
	if yh == 0 {
		return 0
	}
	y := yh.Year()
	if yh.MonthDay() < start {
		y--
	}
	return y
}

package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/MineTakaki/go-utils/conv"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

type (
	//Ymd yyyyMMdd形式で年月日を表す整数型
	Ymd int
)

//ToYmd 年月日からYmd型に変換します
func ToYmd(y, m, d int) (ymd Ymd, err error) {
	_, err = ValidateYmd(y, m, d)
	if err == nil {
		ymd = Ymd(y*10000 + m*100 + d)
	}
	return
}

//ParseYmd Ymd型に変換します
func ParseYmd(i interface{}) (ymd Ymd, err error) {
	err = ymd.Scan(i)
	if err == nil {
		_, err = ymd.Validate()
	}
	return
}

//ParseYmd2 Ymd型に変換します
func ParseYmd2(i interface{}, err *error) (ymd Ymd) {
	var e error
	ymd, e = ParseYmd(i)
	if err != nil {
		*err = e
	}
	return
}

//YmdNow 現在の日付（ローカル）を取得します
func YmdNow() Ymd {
	t := time.Now()
	return Ymd(t.Year()*10000 + int(t.Month())*100 + t.Day())
}

//String string型変換
func (ymd Ymd) String() string {
	if ymd == 0 {
		return ""
	}
	return fmt.Sprintf("%08d", ymd)
}

//CsvFormat CSV出力用のstring型変換
func (ymd Ymd) CsvFormat() string {
	if ymd == 0 {
		return ""
	}
	return strconv.Itoa(int(ymd))
}

//Scan 年月日を読み取ります
func (ymd *Ymd) Scan(i interface{}) (err error) {
	if conv.IsEmpty(i) {
		*ymd = 0
		return
	}

	if tm, ok := i.(time.Time); ok {
		*ymd = Ymd(tm.Year()*10000 + int(tm.Month())*100 + tm.Day())
		return
	}
	if tm, ok := i.(*time.Time); ok {
		*ymd = Ymd(tm.Year()*10000 + int(tm.Month())*100 + tm.Day())
		return
	}
	if n, ok := conv.Int(i); ok {
		*ymd = Ymd(n)
		return
	}
	return ErrValidate
}

//Value driver.Valuerインターフェイスの実装
func (ymd Ymd) Value() (driver.Value, error) {
	if ymd == 0 {
		return nil, nil
	}
	return int64(ymd), nil
}

//UnmarshalJSON json.Unmarshalerインターフェイスの実装
func (ymd *Ymd) UnmarshalJSON(b []byte) (err error) {
	var n json.Number
	if err = json.Unmarshal(b, &n); err != nil {
		err = errors.WithStack(err)
		return
	}
	if n == "" {
		*ymd = 0
		return
	}
	var d int64
	if d, err = strconv.ParseInt(n.String(), 10, 32); err != nil {
		err = errors.WithStack(err)
		return
	}
	x := Ymd(int(d))
	if _, err = x.Validate(); err != nil {
		return
	}
	*ymd = x
	return
}

//MarshalJSON json.Marshalerの実装
func (ymd *Ymd) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(*ymd))
}

//MarshalLogObject zapcore.ObjectMarshalerの実装
func (ymd *Ymd) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("ymd", ymd.String())
	return nil
}

//Validate 年月日が正しいか確認します
func (ymd Ymd) Validate() (bool, error) {
	if ymd == 0 {
		return true, nil
	}
	y, m, d := ymd.Part()
	return ValidateYmd(y, m, d)
}

//Part 年月日の要素を取得します
func (ymd Ymd) Part() (y, m, d int) {
	d = int(ymd) % 100
	x := int(ymd) / 100
	m = x % 100
	y = x / 100
	return
}

//GoTime go言語のTime型に変換します
func (ymd Ymd) GoTime() time.Time {
	y, m, d := ymd.Part()
	y, m, d = AdjustDay(y, m, d)
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
}

//LastDay 最終日を取得します
func LastDay(y, m int) (d int) {
	if m < 1 || m > 12 {
		return
	}
	d = _days[m-1]
	if m == 2 {
		if IsLeapYear(y) {
			d++
		}
	}
	return
}

//Add 年、月、日を加算します（減算はマイナス値を引数にセットします）
func (ymd Ymd) Add(y, m, d int) Ymd {
	year, month, day := ymd.Part()

	//日から計算を行います
	if d != 0 {
		year, month, day = AdjustDay(year, month, day+d)
	}

	if m != 0 {
		year, month = AdjustMonth(year, month+m)
		if maxDay := LastDay(year, month); day > maxDay {
			if d == 0 { //日付の加減算が無い場合は最終日にします
				day = maxDay
			} else if month == 12 { //最終月の場合は翌年
				year++
				month = 1
				day -= maxDay
			} else { //翌月
				month++
				day -= maxDay
			}
		}
	}

	year += y
	return Ymd(year*10000 + month*100 + day)
}

//SetDay 日を指定した値で置き換えます
func (ymd Ymd) SetDay(d int) Ymd {
	xy, xm, _ := ymd.Part()
	if d < 1 {
		d = 1
	} else if lday := LastDay(xy, xm); d > lday {
		d = lday
	}
	return Ymd(xy*10000 + xm*100 + d)
}

//Prev 前日の値を取得します
func (ymd Ymd) Prev() Ymd {
	return ymd.Add(0, 0, -1)
}

//Next 一日後の値を取得します
func (ymd Ymd) Next() Ymd {
	return ymd.Add(0, 0, 1)
}

//Year 年を取得します
func (ymd Ymd) Year() int {
	return int(ymd) / 10000
}

//Month 月を取得します
func (ymd Ymd) Month() int {
	return (int(ymd) / 100) % 100
}

//Day 日を取得します
func (ymd Ymd) Day() int {
	return int(ymd) % 100
}

//YearMonth 年月を取得します
func (ymd Ymd) YearMonth() Ym {
	return Ym(ymd / 100)
}

//MonthDay 月日を取得します
func (ymd Ymd) MonthDay() Md {
	return Md(ymd % 10000)
}

//Between 二つの日付の間に入るか判定します
func (ymd Ymd) Between(f, t Ymd) bool {
	if ymd == 0 || f == 0 || t == 0 {
		return false
	}
	return f <= ymd && ymd <= t
}

//BetweenMonth 2つの月の間に該当するか判定します。m1 > m2の場合は年を跨ぐ範囲として扱います
func (ymd Ymd) BetweenMonth(m1, m2 int) bool {
	if ymd == 0 {
		return false
	}
	return BetweenMonth(ymd.Month(), m1, m2)
}

//BetweenMonthDay 2つの月日の間に該当するか判定します。md1 > md2の場合は年を跨ぐ範囲として扱います
func (ymd Ymd) BetweenMonthDay(md1, md2 Md) bool {
	return ymd.MonthDay().Between(md1, md2)
}

//Min 指定した日付と比較して小さい値を返します
func (ymd Ymd) Min(o Ymd) Ymd {
	if ymd == 0 || ymd > o {
		return o
	}
	return ymd
}

//Max 指定した日付と比較して大きい値を返します
func (ymd Ymd) Max(o Ymd) Ymd {
	if ymd == 0 || ymd < o {
		return o
	}
	return ymd
}

//Days グレゴリウス暦1年1月1日からの経過日数を取得します
func (ymd Ymd) Days() int {
	y, m, d := ymd.Part()

	// 1・2月 → 前年の13・14月
	if m <= 2 {
		y--
		m += 12
	}

	dy := 365 * (y - 1) // 経過年数×365日
	c := y / 100
	dl := (y >> 2) - c + (c >> 2) // うるう年分
	dm := (m*979 - 1033) >> 5     // 1月1日から m 月1日までの日数
	return dy + dl + dm + d - 1
}

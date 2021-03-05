package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"

	"github.com/MineTakaki/go-utils/conv"
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

	year, month, day = AdjustDay(year, month, day)

	//日から計算を行います
	if d > 365 || d < -365 {
		//数値が大きい場合はライブラリで一気に計算します（もしかしたこっちの方が速い？）
		var x time.Month
		year, x, day = ymd.GoTime().AddDate(0, 0, d).Date()
		month = int(x)
	} else {
		year, month, day := ymd.Part()
		day += d
		if month >= 1 && month <= 12 && !((m == 2 && day > 28) || (m == 3 && day <= 0)) {
			if day >= 1 {
				for {
					lday := LastDay(year, month)
					if day <= lday {
						return Ymd((year+y)*10000 + month*100 + day)
					}
					day -= lday
					if month++; month > 12 {
						month = 1
						year++
					}
				}
			}
			for {
				if month--; month < 1 {
					month = 12
					year--
				}
				day += LastDay(year, month)
				if day >= 1 {
					return Ymd((year+y)*10000 + month*100 + day)
				}
			}
		}
	}

	year += y
	month += m
	year += month / 12
	month = ((month-1)%12 + 1)
	if maxDay := _days[month-1]; day > maxDay {
		return Ymd(year*10000 + month*100 + day)
	}
	return Ymd(year*10000 + month*100 + day)
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

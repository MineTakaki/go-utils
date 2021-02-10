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

func tryGetLastDay(m int) (d int) {
	switch m {
	case 1, 3, 5, 7, 8, 10, 12:
		d = 31
	case 4, 6, 9, 11:
		d = 30
	}
	return
}

//Add 年、月、日を加算します（減算はマイナス値を引数にセットします）
func (ymd Ymd) Add(y, m, d int) Ymd {
	//単純な計算で済む場合は計算で済ませる
	if m == 0 {
		year, month, day := ymd.Part()
		year += y
		day += d
		if month >= 1 && month <= 12 && !((m == 2 && day > 28) || (m == 3 && day <= 0)) {
			if day >= 1 {
				for {
					lday := tryGetLastDay(month)
					if lday == 0 { //計算を諦めます
						break
					}
					if day <= lday {
						return Ymd(year*10000 + month*100 + day)
					}
					day -= lday
					if month++; month > 12 {
						month = 1
						year++
					}
				}
			} else {
				for {
					if month--; month < 1 {
						month = 12
						year--
					}
					lday := tryGetLastDay(month)
					if lday == 0 { //計算を諦めます
						break
					}
					day += lday
					if day >= 1 {
						return Ymd(year*10000 + month*100 + day)
					}
				}
			}
		}
	} else if d == 0 {
		year, month, day := ymd.Part()
		year += y
		month += m

		month-- //計算の為に-1します
		year += month / 12
		return Ymd(year*10000 + (month%12+1)*100 + day)
	}
	t := ymd.GoTime().AddDate(y, m, d)
	return Ymd(t.Year()*10000 + int(t.Month())*100 + t.Day())
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
func (ymd Ymd) MonthDay() MdT {
	return MdT(ymd % 10000)
}

//Between 二つの日付の間に入るか判定します
func (ymd Ymd) Between(f, t Ymd) bool {
	return f <= ymd && ymd <= t
}

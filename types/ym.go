package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/MineTakaki/go-utils/conv"
	"github.com/pkg/errors"
)

type (
	//Ym yyyyMM形式で年月を表す整数型
	Ym int

	//MdT MMdd形式で月日を表す整数型
	MdT int
)

//ErrValidate 値が適切でない
var ErrValidate = errors.New("validate error")

//ErrUnkownType 知らない型が指定されました
var ErrUnkownType = errors.New("unkown type")

//ToYm 年月からYm型に変換します
func ToYm(y, m int) (ym Ym, err error) {
	_, err = ValidateYm(y, m)
	if err != nil {
		return
	}
	ym = Ym(y*100 + m)
	return
}

//ParseYm Ym型に変換します
func ParseYm(i interface{}) (ym Ym, err error) {
	err = ym.Scan(i)
	if err == nil {
		_, err = ym.Validate()
	}
	return
}

//ParseYm2 Ym型に変換します
func ParseYm2(i interface{}, err *error) (ym Ym) {
	var e error
	ym, e = ParseYm(i)
	if err != nil {
		*err = e
	}
	return
}

//String string型変換
func (ym Ym) String() string {
	if ym == 0 {
		return ""
	}
	return fmt.Sprintf("%06d", ym)
}

//Validate 年月が正しいか確認します
func (ym Ym) Validate() (bool, error) {
	if ym == 0 {
		return true, nil
	}

	y, m := ym.Part()
	return ValidateYm(y, m)
}

//Year 年を取得します
func (ym Ym) Year() int {
	return int(ym) / 100
}

//Month 月を取得します
func (ym Ym) Month() int {
	return int(ym) % 100
}

//Part 年月の要素を取得します
func (ym Ym) Part() (y, m int) {
	y = int(ym) / 100
	m = int(ym) % 100
	return
}

//Prev 一か月前の値を取得します
func (ym Ym) Prev() Ym {
	return ym.Add(0, -1)
}

//Next 一か月後の値を取得します
func (ym Ym) Next() Ym {
	return ym.Add(0, 1)
}

//Add 年、月を加算します（減算はマイナス値を引数にセットします）
func (ym Ym) Add(dy, dm int) Ym {
	y := ym.Year() + dy
	m := ym.Month() + dm

	for m > 12 {
		y++
		m -= 12
	}
	for m < 1 {
		y--
		m += 12
	}
	return Ym(y*100 + m)
}

//GoTime go言語のTime型に変換します
func (ym Ym) GoTime() time.Time {
	return time.Date(ym.Year(), time.Month(ym.Month()), 1, 0, 0, 0, 0, time.Local)
}

//Term From～To
func (ym Ym) Term() (fm, to Ymd) {
	fm = Ymd(ym*100 + 1)

	m := ym.Month()
	d := tryGetLastDay(m)
	if d != 0 {
		to = Ymd(int(ym)*100 + d)
		return
	}

	//末日がはっきりとしない場合(2月)、ライブラリに任せます
	term := ym.GoTime().AddDate(0, 1, -1)
	to = Ymd(term.Year()*10000 + int(term.Month())*100 + term.Day())

	return
}

func floatToInt(f float64) (n int64, err error) {
	if f < math.MinInt64 || f > math.MaxInt64 {
		err = ErrValidate
	} else {
		n = int64(f)
	}
	return
}

//Scan 年月を読み取ります
func (ym *Ym) Scan(i interface{}) (err error) {
	if conv.IsEmpty(i) {
		*ym = 0
		return nil
	}
	if n, ok := conv.Int(i); ok {
		*ym = Ym(n)
		return
	}
	return ErrValidate
}

//Value driver.Valuerインターフェイスの実装
func (ym Ym) Value() (driver.Value, error) {
	if ym == 0 {
		return nil, nil
	}
	return int64(ym), nil
}

//UnmarshalJSON json.Unmarshalerインターフェイスの実装
func (ym *Ym) UnmarshalJSON(b []byte) (err error) {
	var n int
	err = json.Unmarshal(b, &n)
	if err != nil {
		return
	}
	*ym = Ym(n)
	return
}

//MarshalJSON json.Marshalerの実装
func (ym *Ym) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(*ym))
}

func _validateMD(m, d int) bool {
	if m < 1 || m > 12 {
		return false
	} else if d < 1 || d > 31 {
		return false
	} else if m == 2 {
		if d > 29 {
			return false
		}
	} else if m == 4 || m == 6 || m == 9 || m == 11 {
		if d > 30 {
			return false
		}
	}
	return true
}

//ValidateMd 月日の関連をチェックします（うるう年の考慮はできません）
func ValidateMd(m, d int) (bool, error) {
	if _validateMD(m, d) {
		return true, nil
	}
	return false, errors.Wrapf(ErrValidate, "incorrect month/day value. m:%d, d:%d", m, d)
}

//ValidateYear 年が有効か確認します
func ValidateYear(y int) (bool, error) {
	if y >= 2000 && y <= 2999 {
		return true, nil
	}
	return false, errors.Wrapf(ErrValidate, "%d is not correct as a year value", y)
}

//ValidateMonth 月が有効か確認します
func ValidateMonth(m int) (ok bool, err error) {
	if m < 1 || m > 12 {
		return false, errors.Wrapf(ErrValidate, "%d is not correct as a month value", m)
	}
	return true, nil
}

//ValidateYm 年月が有効か確認します
func ValidateYm(y, m int) (ok bool, err error) {
	ok, err = ValidateYear(y)
	if err != nil {
		return
	}
	ok, err = ValidateMonth(m)
	return
}

//ValidateYmd 年月日が有効か確認します
func ValidateYmd(y, m, d int) (bool, error) {
	if _validateMD(m, d) {
		if y == 9999 { //ターミネータ的な使用目的だけ例外的にOKとする
			if m == 12 && d == 31 {
				return true, nil
			}
		} else if y >= 1998 && y <= 2999 {
			// 2月29日のうるう年以外はOK
			if m != 2 || d != 29 {
				return true, nil
			}

			// 2月29日のうるう年確認を行います
			dt := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
			if dt.Day() == d {
				return true, nil
			}
		}
	}

	return false, errors.Wrapf(ErrValidate, "incorrect date value. y:%d, m:%d, d:%d", y, m, d)
}

//String string型変換
func (md MdT) String() string {
	if md == 0 {
		return ""
	}
	return fmt.Sprintf("%04d", md)
}

//CsvFormat CSV出力用のstring型変換
func (ym Ym) CsvFormat() string {
	if ym == 0 {
		return ""
	}
	return strconv.Itoa(int(ym))
}

//Scan 文字列から月日を読み取ります
func (md *MdT) Scan(s string) (err error) {
	if s == "" {
		*md = 0
		return
	}
	var tm time.Time
	tm, err = time.Parse("0102", s) // MMDD形式
	if err == nil {
		*md = MdT(int(tm.Month())*100 + tm.Day())
	}
	return
}

//Between 二つの日付の間に入るか判定します
func (md MdT) Between(f, t MdT) bool {
	if md == 0 || f == 0 || t == 0 {
		return false
	}
	if f > t {
		return md >= f || md <= t
	}
	return f <= md && md >= t
}

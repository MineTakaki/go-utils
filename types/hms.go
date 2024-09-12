package types

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/MineTakaki/go-utils/conv"
	"github.com/MineTakaki/go-utils/errors"
	"go.uber.org/zap/zapcore"
)

type (
	//Hms hhmmss形式で時分秒を表す整数型
	Hms int

	//HmsSlice Hms型のスライス
	HmsSlice []Hms
)

func (hms HmsSlice) Len() int           { return len(hms) }
func (hms HmsSlice) Less(i, j int) bool { return hms[i] < hms[j] }
func (hms HmsSlice) Swap(i, j int)      { hms[i], hms[j] = hms[j], hms[i] }

// ToHms 時分秒からHms型に変換します
func ToHms(h, m, s int) (hms Hms, err error) {
	_, err = ValidateHms(h, m, s)
	if err != nil {
		return
	}
	hms = Hms(h*10000 + m*100 + s)
	return
}

// ParseHms Hms型に変換します
func ParseHms(i interface{}) (hms Hms, err error) {
	err = hms.Scan(i)
	if err == nil {
		_, err = hms.Validate()
	}
	return
}

// ParseHms2 Hms型に変換します
func ParseHms2(i interface{}, err *error) (hms Hms) {
	var e error
	hms, e = ParseHms(i)
	if err != nil {
		*err = e
	}
	return
}

// String string型変換
func (hms Hms) String() string {
	return ItoaZeroFilled(int(hms), 6)
}

// FormatHms hhmmss形式でstring型に整形して変換します
func (hms Hms) FormatHms(sep string, zeroSuppress bool) string {
	if hms == 0 {
		return ""
	}
	h, m, s := hms.Part()
	sb := strings.Builder{}
	sb.Grow(len(sep)*2 + 6)
	if zeroSuppress {
		sb.WriteString(strconv.Itoa(h))
		sb.WriteString(sep)
		sb.WriteString(strconv.Itoa(m))
		sb.WriteString(sep)
		sb.WriteString(strconv.Itoa(s))
		return sb.String()
	}
	sb.WriteString(fillZero2(h))
	sb.WriteString(sep)
	sb.WriteString(fillZero2(m))
	sb.WriteString(sep)
	sb.WriteString(fillZero2(s))
	return sb.String()
}

// Value driver.Valuerインターフェイスの実装
func (hms Hms) Value() (driver.Value, error) {
	if hms == 0 {
		return nil, nil
	}
	return int64(hms), nil
}

// UnmarshalJSON json.Unmarshalerインターフェイスの実装
func (hms *Hms) UnmarshalJSON(b []byte) (err error) {
	var s interface{}
	if err = errors.WithStack(json.Unmarshal(b, &s)); err != nil {
		return
	}
	var x Hms
	if x, err = ParseHms(s); err != nil {
		return
	}
	*hms = x
	return
}

// MarshalJSON json.Marshalerの実装
func (hms *Hms) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(*hms))
}

// MarshalLogObject zapcore.ObjectMarshalerの実装
func (hms *Hms) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("hms", hms.String())
	return nil
}

// Validate 年月が正しいか確認します
func (hms Hms) Validate() (bool, error) {
	if hms == 0 {
		return true, nil
	}

	h, m, s := hms.Part()
	return ValidateHms(h, m, s)
}

// Hour 月を取得します
func (hms Hms) Hour() int {
	return int(hms) / 10000
}

// Minute 月を取得します
func (hms Hms) Minute() int {
	return (abs(int(hms)) / 100) % 100
}

// Second 秒を取得します
func (hms Hms) Second() int {
	return abs(int(hms)) % 100
}

// Part 時分秒の要素を取得します
func (hms Hms) Part() (h, m, s int) {
	s = abs(int(hms)) % 100
	n := int(hms) / 100
	m = abs(n) % 100
	h = n / 100
	return
}

// Parts 時分秒の要素を配列で取得します
func (hms Hms) Parts() []int {
	v := make([]int, 3)
	v[0], v[1], v[2] = hms.Part()
	return v
}

// Prev 1秒前の時間を取得します
func (hms Hms) Prev() Hms {
	return hms.Add(0, 0, -1)
}

// Next 1秒後の値を取得します
func (hms Hms) Next() Hms {
	return hms.Add(0, 0, 1)
}

// AdjustHms 時分秒の加減算後の正しい時間を取得します
func AdjustHms(h, m, s int) (qh, qm, qs int) {
	n := (h*60+m)*60 + s
	qs = abs(n) % 60
	n /= 60
	qm = abs(n) % 60
	qh = n / 60
	return
}

// Add 時、分、秒を加算します（減算はマイナス値を引数にセットします）
func (hms Hms) Add(dh, dm, ds int) Hms {
	if hms == 0 {
		return 0
	}
	h, m, s := hms.Part()

	//状態を正常化します
	h, m, s = AdjustHms(h, m, s+ds)

	return Hms(h*10000 + m*100 + s)
}

// Adjust 時分秒を正しい形式に訂正します
func (hms Hms) Adjust() Hms {
	return hms.Add(0, 0, 0)
}

// Scan 文字列から時分秒を読み取ります
func (hms *Hms) Scan(i interface{}) (err error) {
	if conv.IsEmpty(i) {
		*hms = 0
		return nil
	}
	if n, ok := conv.Int(i); ok {
		*hms = Hms(n)
		return
	}
	return errors.WithStack(ErrValidate)
}

// Between 二つの日付の間に入るか判定します
func (hms Hms) Between(f, t Hms) bool {
	if hms == 0 || f == 0 || t == 0 {
		return false
	}
	if f > t {
		return hms >= f || hms <= t
	}
	return f <= hms && hms >= t
}

// Compare Hms同志を比較します
func (hms Hms) Compare(o Hms) int {
	if hms < o {
		return -1
	}
	if hms > o {
		return 1
	}
	return 0
}

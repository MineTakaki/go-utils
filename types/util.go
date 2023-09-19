package types

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// ZeroPrefix 指定文字列をwidthで指定した桁で0埋めしてstringに変換します
//
//	半角数字を想定しているのでマルチバイト文字列を渡すのは禁止です
func ZeroPrefix(s string, width int) string {
	switch {
	default:
		sb := strings.Builder{}
		sb.Grow(6)
		for i := width - len(s); i >= 0; i-- {
			sb.WriteByte('0')
		}
		sb.WriteString(s)
		return sb.String()
	case len(s) == width:
		return s
	case len(s) > width:
		return s[len(s)-width:]
	}
}

// ItoaZeroFilled intの値をwidthで指定した桁で0埋めしてstringに変換します
//
//	widthからあふれた場合は上位の桁を切り捨てます
func ItoaZeroFilled(n, width int) string {
	return ZeroPrefix(strconv.Itoa(n), width)
}

// ItoaZeroFilled64 int64の値をwidthで指定した桁で0埋めしてstringに変換します
//
//	widthからあふれた場合は上位の桁を切り捨てます
func ItoaZeroFilled64(n int64, width int) string {
	return ZeroPrefix(strconv.FormatInt(n, 10), width)
}

func fillZero4(y int) string {
	if y < 0 {
		y = -y
	}
	if y < 10 {
		return "000" + strconv.Itoa(y)
	} else if y < 100 {
		return "00" + strconv.Itoa(y)
	} else if y < 1000 {
		return "0" + strconv.Itoa(y)
	}
	return strconv.Itoa(y)
}

func fillZero2(n int) string {
	if n < 0 {
		n = -n
	}
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}
func _validateMD(m, d int) bool {
	if m < 1 || m > 12 {
		return false
	} else if d < 1 || d > 31 {
		return false
	} else if m == 2 {
		if d > 28 {
			return false
		}
	} else if m == 4 || m == 6 || m == 9 || m == 11 {
		if d > 30 {
			return false
		}
	}
	return true
}

// ValidateMd 月日の関連をチェックします（うるう年の考慮はできません）
func ValidateMd(m, d int) (bool, error) {
	if _validateMD(m, d) {
		return true, nil
	}
	return false, errors.Wrapf(ErrValidate, "incorrect month/day value. m:%d, d:%d", m, d)
}

// ValidateYear 年が有効か確認します
func ValidateYear(y int) (bool, error) {
	if y >= 1998 && y <= 2999 {
		return true, nil
	}
	return false, errors.Wrapf(ErrValidate, "%d is not correct as a year value", y)
}

// ValidateMonth 月が有効か確認します
func ValidateMonth(m int) (ok bool, err error) {
	if m < 1 || m > 12 {
		return false, errors.Wrapf(ErrValidate, "%d is not correct as a month value", m)
	}
	return true, nil
}

// ValidateYm 年月が有効か確認します
func ValidateYm(y, m int) (ok bool, err error) {
	ok, err = ValidateYear(y)
	if err != nil {
		return
	}
	ok, err = ValidateMonth(m)
	return
}

// ValidateYmd 年月日が有効か確認します
func ValidateYmd(y, m, d int) (bool, error) {
	if y == 9999 { //ターミネータ的な使用目的だけ例外的にOKとする
		if m == 12 && d == 31 {
			return true, nil
		} else if m == 99 && d == 99 {
			return true, nil
		}
	} else if y >= 1998 && y <= 2999 && m >= 1 && m <= 12 && d >= 1 {
		if lday := LastDay(y, m); d <= lday {
			return true, nil
		}
	}
	return false, errors.Wrapf(ErrValidate, "incorrect date value. y:%d, m:%d, d:%d", y, m, d)
}

// ValidateHms 時分秒が有効か確認します
func ValidateHms(h, m, s int) (bool, error) {
	if h >= 0 && h <= 23 && m >= 0 && m <= 59 && s >= 0 && s <= 59 {
		return true, nil
	}
	return false, errors.Wrapf(ErrValidate, "incorrect time value. h:%d, m:%d, s:%d", h, m, s)
}

// ValidateYmdhms 年月日時分秒が有効か確認します
func ValidateYmdhms(y, m, d, h, n, s int) (bool, error) {
	if y >= 1998 && y <= 2999 && m >= 1 && m <= 12 && d >= 1 && d <= 31 {
		if h >= 0 && h <= 23 && n >= 0 && n <= 59 && s >= 0 && s <= 59 {
			if lday := LastDay(y, m); d <= lday {
				return true, nil
			}
		}
	}
	return false, errors.Wrapf(ErrValidate, "incorrect date value. y:%d, m:%d, d:%d, h:%d, n:%d, s:%d", y, m, d, h, n, s)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func abs64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

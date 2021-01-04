package decimal

import (
	"database/sql/driver"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/zap/zapcore"
)

type (
	//Decimal Decimal型
	// github.com/shopspring/decimal.Decimal
	Decimal struct {
		decimal.Decimal
	}
)

// DivisionPrecision is the number of decimal places in the result when it
// doesn't divide exactly.
//
// Example:
//
//     d1 := decimal.NewFromFloat(2).Div(decimal.NewFromFloat(3))
//     d1.String() // output: "0.6666666666666667"
//     d2 := decimal.NewFromFloat(2).Div(decimal.NewFromFloat(30000))
//     d2.String() // output: "0.0000666666666667"
//     d3 := decimal.NewFromFloat(20000).Div(decimal.NewFromFloat(3))
//     d3.String() // output: "6666.6666666666666667"
//     decimal.DivisionPrecision = 3
//     d4 := decimal.NewFromFloat(2).Div(decimal.NewFromFloat(3))
//     d4.String() // output: "0.667"
//
var DivisionPrecision = 16

// MarshalJSONWithoutQuotes should be set to true if you want the decimal to
// be JSON marshaled as a number, instead of as a string.
// WARNING: this is dangerous for decimals with many digits, since many JSON
// unmarshallers (ex: Javascript's) will unmarshal JSON numbers to IEEE 754
// double-precision floating point numbers, which means you can potentially
// silently lose precision.
var MarshalJSONWithoutQuotes = false

// Zero constant, to make computations faster.
// Zero should never be compared with == or != directly, please use decimal.Equal or decimal.Cmp instead.
var Zero = New(0, 1)

// Cent constant, to make computations faster.
var Cent = New(100, 1)

//ErrScan Scan error
var ErrScan = errors.New("scan value error")

// Scan implements the sql.Scanner interface for database deserialization.
func (d *Decimal) Scan(value interface{}) error {
	if x, ok := ValueOf(value); ok {
		*d = x
		return nil
	}
	return errors.WithStack(ErrScan)
}

//MarshalLogObject implements of zapcore.ObjectMarshaler interface.
func (d Decimal) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("decimal", d.String())
	return nil
}

// Value implements the driver.Valuer interface for database serialization.
func (d Decimal) Value() (driver.Value, error) {
	return d.Decimal.Value()
}

//Nullable Decimal型をNullDecimal型に変換します
func (d Decimal) Nullable() NullDecimal {
	return NullDecimal{
		Decimal: d,
		Valid:   true,
	}
}

// Abs returns the absolute value of the decimal.
func (d Decimal) Abs() Decimal {
	return Decimal{d.Decimal.Abs()}
}

// Add returns d + d2.
func (d Decimal) Add(d2 Decimal) Decimal {
	return Decimal{d.Decimal.Add(d2.Decimal)}
}

// Atan returns the arctangent, in radians, of x.
func (d Decimal) Atan() Decimal {
	return Decimal{d.Decimal.Atan()}
}

// Ceil returns the nearest integer value greater than or equal to d.
func (d Decimal) Ceil() Decimal {
	return Decimal{d.Decimal.Ceil()}
}

// Cos returns the cosine of the radian argument x.
func (d Decimal) Cos() Decimal {
	return Decimal{d.Decimal.Cos()}
}

// Div returns d / d2. If it doesn't divide exactly, the result will have
// DivisionPrecision digits after the decimal point.
func (d Decimal) Div(d2 Decimal) Decimal {
	return Decimal{d.Decimal.DivRound(d2.Decimal, int32(DivisionPrecision))}
}

// DivRound divides and rounds to a given precision
// i.e. to an integer multiple of 10^(-precision)
//   for a positive quotient digit 5 is rounded up, away from 0
//   if the quotient is negative then digit 5 is rounded down, away from 0
// Note that precision<0 is allowed as input.
func (d Decimal) DivRound(d2 Decimal, precision int32) Decimal {
	return Decimal{d.Decimal.DivRound(d2.Decimal, precision)}
}

// Floor returns the nearest integer value less than or equal to d.
func (d Decimal) Floor() Decimal {
	return Decimal{d.Decimal.Floor()}
}

// Mod returns d % d2.
func (d Decimal) Mod(d2 Decimal) Decimal {
	return Decimal{d.Decimal.Mod(d2.Decimal)}
}

// Mul returns d * d2.
func (d Decimal) Mul(d2 Decimal) Decimal {
	return Decimal{d.Decimal.Mul(d2.Decimal)}
}

// Neg returns -d.
func (d Decimal) Neg() Decimal {
	return Decimal{d.Decimal.Neg()}
}

// Pow returns d to the power d2
func (d Decimal) Pow(d2 Decimal) Decimal {
	return Decimal{d.Decimal.Pow(d2.Decimal)}
}

// Cmp compares the numbers represented by d and d2 and returns:
//
//     -1 if d <  d2
//      0 if d == d2
//     +1 if d >  d2
//
func (d Decimal) Cmp(d2 Decimal) int {
	return d.Decimal.Cmp(d2.Decimal)
}

// Equal returns whether the numbers represented by d and d2 are equal.
func (d Decimal) Equal(d2 Decimal) bool {
	return d.Decimal.Equal(d2.Decimal)
}

// Equals is deprecated, please use Equal method instead
func (d Decimal) Equals(d2 Decimal) bool {
	return d.Decimal.Equals(d2.Decimal)
}

// GreaterThan (GT) returns true when d is greater than d2.
func (d Decimal) GreaterThan(d2 Decimal) bool {
	return d.Decimal.GreaterThan(d2.Decimal)
}

// GreaterThanOrEqual (GTE) returns true when d is greater than or equal to d2.
func (d Decimal) GreaterThanOrEqual(d2 Decimal) bool {
	return d.Decimal.GreaterThanOrEqual(d2.Decimal)
}

// LessThan (LT) returns true when d is less than d2.
func (d Decimal) LessThan(d2 Decimal) bool {
	return d.Decimal.LessThan(d2.Decimal)
}

// LessThanOrEqual (LTE) returns true when d is less than or equal to d2.
func (d Decimal) LessThanOrEqual(d2 Decimal) bool {
	return d.Decimal.LessThanOrEqual(d2.Decimal)
}

// Sign returns:
//
//	-1 if d <  0
//	 0 if d == 0
//	+1 if d >  0
//
func (d Decimal) Sign() int {
	return d.Decimal.Sign()
}

// IsPositive return
//
//	true if d > 0
//	false if d == 0
//	false if d < 0
func (d Decimal) IsPositive() bool {
	return d.Decimal.Sign() == 1
}

// IsNegative return
//
//	true if d < 0
//	false if d == 0
//	false if d > 0
func (d Decimal) IsNegative() bool {
	return d.Decimal.Sign() == -1
}

// IsZero return
//
//	true if d == 0
//	false if d > 0
//	false if d < 0
func (d Decimal) IsZero() bool {
	return d.Decimal.IsZero()
}

// Exponent returns the exponent, or scale component of the decimal.
func (d Decimal) Exponent() int32 {
	return d.Decimal.Exponent()
}

// Coefficient returns the coefficient of the decimal.  It is scaled by 10^Exponent()
func (d Decimal) Coefficient() *big.Int {
	return d.Decimal.Coefficient()
}

// IntPart returns the integer component of the decimal.
func (d Decimal) IntPart() int64 {
	return d.Decimal.IntPart()
}

// BigInt returns integer component of the decimal as a BigInt.
func (d Decimal) BigInt() *big.Int {
	return d.Decimal.BigInt()
}

// BigFloat returns decimal as BigFloat.
// Be aware that casting decimal to BigFloat might cause a loss of precision.
func (d Decimal) BigFloat() *big.Float {
	return d.Decimal.BigFloat()
}

// Rat returns a rational number representation of the decimal.
func (d Decimal) Rat() *big.Rat {
	return d.Decimal.Rat()
}

// Float64 returns the nearest float64 value for d and a bool indicating
// whether f represents d exactly.
// For more details, see the documentation for big.Rat.Float64
func (d Decimal) Float64() (f float64, exact bool) {
	return d.Decimal.Float64()
}

// String returns the string representation of the decimal
// with the fixed point.
//
// Example:
//
//     d := New(-12345, -3)
//     println(d.String())
//
// Output:
//
//     -12.345
//
func (d Decimal) String() string {
	return d.Decimal.String()
}

// StringFixed returns a rounded fixed-point string with places digits after
// the decimal point.
//
// Example:
//
// 	   NewFromFloat(0).StringFixed(2) // output: "0.00"
// 	   NewFromFloat(0).StringFixed(0) // output: "0"
// 	   NewFromFloat(5.45).StringFixed(0) // output: "5"
// 	   NewFromFloat(5.45).StringFixed(1) // output: "5.5"
// 	   NewFromFloat(5.45).StringFixed(2) // output: "5.45"
// 	   NewFromFloat(5.45).StringFixed(3) // output: "5.450"
// 	   NewFromFloat(545).StringFixed(-1) // output: "550"
//
func (d Decimal) StringFixed(places int32) string {
	return d.Decimal.StringFixed(places)
}

// StringFixedBank returns a banker rounded fixed-point string with places digits
// after the decimal point.
//
// Example:
//
// 	   NewFromFloat(0).StringFixedBank(2) // output: "0.00"
// 	   NewFromFloat(0).StringFixedBank(0) // output: "0"
// 	   NewFromFloat(5.45).StringFixedBank(0) // output: "5"
// 	   NewFromFloat(5.45).StringFixedBank(1) // output: "5.4"
// 	   NewFromFloat(5.45).StringFixedBank(2) // output: "5.45"
// 	   NewFromFloat(5.45).StringFixedBank(3) // output: "5.450"
// 	   NewFromFloat(545).StringFixedBank(-1) // output: "540"
//
func (d Decimal) StringFixedBank(places int32) string {
	return d.Decimal.StringFixedBank(places)
}

// StringFixedCash returns a Swedish/Cash rounded fixed-point string. For
// more details see the documentation at function RoundCash.
func (d Decimal) StringFixedCash(interval uint8) string {
	return d.Decimal.StringFixedCash(interval)
}

// QuoRem does divsion with remainder
// d.QuoRem(d2,precision) returns quotient q and remainder r such that
//   d = d2 * q + r, q an integer multiple of 10^(-precision)
//   0 <= r < abs(d2) * 10 ^(-precision) if d>=0
//   0 >= r > -abs(d2) * 10 ^(-precision) if d<0
// Note that precision<0 is allowed as input.
func (d Decimal) QuoRem(d2 Decimal, precision int32) (Decimal, Decimal) {
	q1, q2 := d.Decimal.QuoRem(d2.Decimal, precision)
	return Decimal{q1}, Decimal{q2}
}

// Round rounds the decimal to places decimal places.
// If places < 0, it will round the integer part to the nearest 10^(-places).
//
// Example:
//
// 	   NewFromFloat(5.45).Round(1).String() // output: "5.5"
// 	   NewFromFloat(545).Round(-1).String() // output: "550"
//
func (d Decimal) Round(places int32) Decimal {
	return Decimal{d.Decimal.Round(places)}
}

// RoundBank rounds the decimal to places decimal places.
// If the final digit to round is equidistant from the nearest two integers the
// rounded value is taken as the even number
//
// If places < 0, it will round the integer part to the nearest 10^(-places).
//
// Examples:
//
// 	   NewFromFloat(5.45).Round(1).String() // output: "5.4"
// 	   NewFromFloat(545).Round(-1).String() // output: "540"
// 	   NewFromFloat(5.46).Round(1).String() // output: "5.5"
// 	   NewFromFloat(546).Round(-1).String() // output: "550"
// 	   NewFromFloat(5.55).Round(1).String() // output: "5.6"
// 	   NewFromFloat(555).Round(-1).String() // output: "560"
//
func (d Decimal) RoundBank(places int32) Decimal {
	return Decimal{d.Decimal.RoundBank(places)}
}

// RoundCash aka Cash/Penny/öre rounding rounds decimal to a specific
// interval. The amount payable for a cash transaction is rounded to the nearest
// multiple of the minimum currency unit available. The following intervals are
// available: 5, 10, 25, 50 and 100; any other number throws a panic.
//	    5:   5 cent rounding 3.43 => 3.45
// 	   10:  10 cent rounding 3.45 => 3.50 (5 gets rounded up)
// 	   25:  25 cent rounding 3.41 => 3.50
// 	   50:  50 cent rounding 3.75 => 4.00
// 	  100: 100 cent rounding 3.50 => 4.00
// For more details: https://en.wikipedia.org/wiki/Cash_rounding
func (d Decimal) RoundCash(interval uint8) Decimal {
	return Decimal{d.Decimal.RoundCash(interval)}
}

// Shift shifts the decimal in base 10.
// It shifts left when shift is positive and right if shift is negative.
// In simpler terms, the given value for shift is added to the exponent
// of the decimal.
func (d Decimal) Shift(shift int32) Decimal {
	return Decimal{d.Decimal.Shift(shift)}
}

// Sin returns the sine of the radian argument x.
func (d Decimal) Sin() Decimal {
	return Decimal{d.Decimal.Sin()}
}

// Sub returns d - d2.
func (d Decimal) Sub(d2 Decimal) Decimal {
	return Decimal{d.Decimal.Sub(d2.Decimal)}
}

// Tan returns the tangent of the radian argument x.
func (d Decimal) Tan() Decimal {
	return Decimal{d.Decimal.Tan()}
}

// Truncate truncates off digits from the number, without rounding.
//
// NOTE: precision is the last digit that will not be truncated (must be >= 0).
//
// Example:
//
//     decimal.NewFromString("123.456").Truncate(2).String() // "123.45"
//
func (d Decimal) Truncate(precision int32) Decimal {
	return Decimal{d.Decimal.Truncate(precision)}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Decimal) UnmarshalJSON(decimalBytes []byte) error {
	if string(decimalBytes) == "null" {
		return nil
	}

	str := unquoteIfQuoted(decimalBytes)

	decimal, err := NewFromString(str)
	*d = decimal
	if err != nil {
		return fmt.Errorf("error decoding string '%s': %s", str, err)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (d Decimal) MarshalJSON() ([]byte, error) {
	var str string
	if MarshalJSONWithoutQuotes {
		str = d.String()
	} else {
		str = "\"" + d.String() + "\""
	}
	return []byte(str), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface. As a string representation
// is already used when encoding to text, this method stores that string as []byte
func (d *Decimal) UnmarshalBinary(data []byte) error {
	return d.Decimal.UnmarshalBinary(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (d Decimal) MarshalBinary() (data []byte, err error) {
	return d.Decimal.MarshalBinary()
}

// New returns a new fixed-point decimal, value * 10 ^ exp.
func New(value int64, exp int32) Decimal {
	return Decimal{decimal.New(value, exp)}
}

// NewFromInt converts a int64 to Decimal.
//
// Example:
//
//     NewFromInt(123).String() // output: "123"
//     NewFromInt(-10).String() // output: "-10"
func NewFromInt(value int64) Decimal {
	return Decimal{decimal.NewFromInt(value)}
}

// NewFromInt32 converts a int32 to Decimal.
//
// Example:
//
//     NewFromInt(123).String() // output: "123"
//     NewFromInt(-10).String() // output: "-10"
func NewFromInt32(value int32) Decimal {
	return Decimal{decimal.NewFromInt32(value)}
}

// NewFromBigInt returns a new Decimal from a big.Int, value * 10 ^ exp
func NewFromBigInt(value *big.Int, exp int32) Decimal {
	return Decimal{decimal.NewFromBigInt(value, exp)}
}

// NewFromString returns a new Decimal from a string representation.
// Trailing zeroes are not trimmed.
//
// Example:
//
//     d, err := NewFromString("-123.45")
//     d2, err := NewFromString(".0001")
//     d3, err := NewFromString("1.47000")
//
func NewFromString(value string) (Decimal, error) {
	d, err := decimal.NewFromString(value)
	return Decimal{d}, err
}

// RequireFromString returns a new Decimal from a string representation
// or panics if NewFromString would have returned an error.
//
// Example:
//
//     d := RequireFromString("-123.45")
//     d2 := RequireFromString(".0001")
//
func RequireFromString(value string) Decimal {
	return Decimal{decimal.RequireFromString(value)}
}

// NewFromFloat converts a float64 to Decimal.
//
// The converted number will contain the number of significant digits that can be
// represented in a float with reliable roundtrip.
// This is typically 15 digits, but may be more in some cases.
// See https://www.exploringbinary.com/decimal-precision-of-binary-floating-point-numbers/ for more information.
//
// For slightly faster conversion, use NewFromFloatWithExponent where you can specify the precision in absolute terms.
//
// NOTE: this will panic on NaN, +/-inf
func NewFromFloat(value float64) Decimal {
	return Decimal{decimal.NewFromFloat(value)}
}

// NewFromFloat32 converts a float32 to Decimal.
//
// The converted number will contain the number of significant digits that can be
// represented in a float with reliable roundtrip.
// This is typically 6-8 digits depending on the input.
// See https://www.exploringbinary.com/decimal-precision-of-binary-floating-point-numbers/ for more information.
//
// For slightly faster conversion, use NewFromFloatWithExponent where you can specify the precision in absolute terms.
//
// NOTE: this will panic on NaN, +/-inf
func NewFromFloat32(value float32) Decimal {
	return Decimal{decimal.NewFromFloat32(value)}
}

// NewFromFloatWithExponent converts a float64 to Decimal, with an arbitrary
// number of fractional digits.
//
// Example:
//
//     NewFromFloatWithExponent(123.456, -2).String() // output: "123.46"
//
func NewFromFloatWithExponent(value float64, exp int32) Decimal {
	return Decimal{decimal.NewFromFloatWithExponent(value, exp)}
}

// Min returns the smallest Decimal that was passed in the arguments.
//
// To call this function with an array, you must do:
//
//     Min(arr[0], arr[1:]...)
//
// This makes it harder to accidentally call Min with 0 arguments.
func Min(first Decimal, rest ...Decimal) Decimal {
	ans := first
	for _, item := range rest {
		if item.Cmp(ans) < 0 {
			ans = item
		}
	}
	return ans
}

// Max returns the largest Decimal that was passed in the arguments.
//
// To call this function with an array, you must do:
//
//     Max(arr[0], arr[1:]...)
//
// This makes it harder to accidentally call Max with 0 arguments.
func Max(first Decimal, rest ...Decimal) Decimal {
	ans := first
	for _, item := range rest {
		if item.Cmp(ans) > 0 {
			ans = item
		}
	}
	return ans
}

// Sum returns the combined total of the provided first and rest Decimals
func Sum(first Decimal, rest ...Decimal) Decimal {
	total := first
	for _, item := range rest {
		total = total.Add(item)
	}

	return total
}

// Avg returns the average value of the provided first and rest Decimals
func Avg(first Decimal, rest ...Decimal) Decimal {
	count := New(int64(len(rest)+1), 0)
	sum := Sum(first, rest...)
	return sum.Div(count)
}

// RescalePair rescales two decimals to common exponential value (minimal exp of both decimals)
func RescalePair(d1 Decimal, d2 Decimal) (Decimal, Decimal) {
	p1, p2 := decimal.RescalePair(d1.Decimal, d2.Decimal)
	return Decimal{p1}, Decimal{p2}
}

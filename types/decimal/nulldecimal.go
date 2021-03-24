package decimal

import (
	"database/sql"
	"database/sql/driver"

	"go.uber.org/zap/zapcore"
)

type (
	//NullDecimal NullDecimal
	NullDecimal struct {
		Decimal Decimal
		Valid   bool
	}
)

var _ driver.Valuer = &NullDecimal{}

//Null NullDecimal型のNull値
var Null = NullDecimal{}

//NullZero NullDecimal型のZero値
var NullZero = NullDecimal{Valid: true}

//NullCent NullDecimal型の100
var NullCent = NullDecimal{Decimal: Cent, Valid: true}

//String string return
func (d NullDecimal) String() string {
	if !d.Valid {
		return ""
	}
	return d.Decimal.String()
}

//Equal NullDecimal同士が同じか確認します
func (d NullDecimal) Equal(o NullDecimal) bool {
	if d.Valid {
		if !o.Valid {
			return false
		}
		return d.Decimal.Equal(o.Decimal)
	}
	if o.Valid {
		return false
	}
	return true //両方ともにNULL
}

//EqualNZ NullDecimal同士が同じか確認します(NULLはZEROと判断します)
func (d NullDecimal) EqualNZ(o NullDecimal) bool {
	if d.Valid {
		if !o.Valid {
			return d.Decimal.IsZero()
		}
		return d.Decimal.Equal(o.Decimal)
	}
	if o.Valid {
		return o.Decimal.IsZero()
	}
	return true //両方ともにNULL
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *NullDecimal) Scan(value interface{}) error {
	if value == nil {
		d.Decimal = Zero
		d.Valid = false
		return nil
	}
	d.Valid = true
	return d.Decimal.Scan(value)
}

//MarshalLogObject implements of zapcore.ObjectMarshaler interface.
func (d NullDecimal) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("decimal", d.Decimal.String())
	enc.AddBool("valid", d.Valid)
	return nil
}

// Value implements the driver.Valuer interface for database serialization.
func (d NullDecimal) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Decimal.Value()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *NullDecimal) UnmarshalJSON(decimalBytes []byte) error {
	if string(decimalBytes) == "null" {
		d.Valid = false
		return nil
	}
	d.Valid = true
	return d.Decimal.UnmarshalJSON(decimalBytes)
}

// MarshalJSON implements the json.Marshaler interface.
func (d NullDecimal) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}
	return d.Decimal.MarshalJSON()
}

// Abs returns the absolute value of the decimal.
func (d NullDecimal) Abs() NullDecimal {
	if !d.Valid {
		return Null
	}
	return d.Decimal.Abs().Nullable()
}

// Add returns d + d2.
func (d NullDecimal) Add(d2 interface{}) NullDecimal {
	if v, ok := ValueOf(d2); ok {
		if d.Valid {
			return d.Decimal.Add(v).Nullable()
		}
		return v.Nullable()
	}
	return d
}

// Atan returns the arctangent, in radians, of x.
func (d NullDecimal) Atan() NullDecimal {
	if !d.Valid {
		return Null
	}
	return d.Decimal.Atan().Nullable()
}

// Ceil returns the nearest integer value greater than or equal to d.
func (d NullDecimal) Ceil() NullDecimal {
	if d.Valid {
		return Null
	}
	return d.Decimal.Ceil().Nullable()
}

// Cos returns the cosine of the radian argument x.
func (d NullDecimal) Cos() NullDecimal {
	if d.Valid {
		return Null
	}
	return d.Decimal.Cos().Nullable()
}

// Div returns d / d2. If it doesn't divide exactly, the result will have
// DivisionPrecision digits after the decimal point.
func (d NullDecimal) Div(d2 interface{}) NullDecimal {
	return d.DivRound(d2, int32(DivisionPrecision))
}

// DivRound divides and rounds to a given precision
// i.e. to an integer multiple of 10^(-precision)
//   for a positive quotient digit 5 is rounded up, away from 0
//   if the quotient is negative then digit 5 is rounded down, away from 0
// Note that precision<0 is allowed as input.
func (d NullDecimal) DivRound(d2 interface{}, precision int32) NullDecimal {
	if !d.Valid {
		return Null
	}
	if v, ok := ValueOf(d2); ok {
		return d.Decimal.DivRound(v, precision).Nullable()
	}
	return d
}

// Floor returns the nearest integer value less than or equal to d.
func (d NullDecimal) Floor() NullDecimal {
	if !d.Valid {
		return Null
	}
	return d.Decimal.Floor().Nullable()
}

// Mod returns d % d2.
func (d NullDecimal) Mod(d2 interface{}) NullDecimal {
	if !d.Valid {
		return Null
	}
	if v, ok := ValueOf(d2); ok {
		return d.Decimal.Mod(v).Nullable()
	}
	return Null
}

// Mul returns d * d2.
func (d NullDecimal) Mul(d2 interface{}) NullDecimal {
	if !d.Valid {
		return Null
	}
	if v, ok := ValueOf(d2); ok {
		return d.Decimal.Mul(v).Nullable()
	}
	return Null
}

// Neg returns -d.
func (d NullDecimal) Neg() NullDecimal {
	if !d.Valid {
		return Null
	}
	return d.Decimal.Neg().Nullable()
}

// Pow returns d to the power d2
func (d NullDecimal) Pow(d2 interface{}) NullDecimal {
	if !d.Valid {
		return Null
	}
	if v, ok := ValueOf(d2); ok {
		return d.Decimal.Pow(v).Nullable()
	}
	return Null
}

// QuoRem does divsion with remainder
// d.QuoRem(d2,precision) returns quotient q and remainder r such that
//   d = d2 * q + r, q an integer multiple of 10^(-precision)
//   0 <= r < abs(d2) * 10 ^(-precision) if d>=0
//   0 >= r > -abs(d2) * 10 ^(-precision) if d<0
// Note that precision<0 is allowed as input.
func (d NullDecimal) QuoRem(d2 interface{}, precision int32) (NullDecimal, NullDecimal) {
	if d.Valid {
		if v, ok := ValueOf(d2); ok {
			q1, q2 := d.Decimal.QuoRem(v, precision)
			return q1.Nullable(), q2.Nullable()
		}
	}
	return Null, Null
}

// Round rounds the decimal to places decimal places.
// If places < 0, it will round the integer part to the nearest 10^(-places).
//
// Example:
//
// 	   NewFromFloat(5.45).Round(1).String() // output: "5.5"
// 	   NewFromFloat(545).Round(-1).String() // output: "550"
//
func (d NullDecimal) Round(places int32) NullDecimal {
	if !d.Valid {
		return Null
	}
	return d.Decimal.Round(places).Nullable()
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
func (d NullDecimal) RoundBank(places int32) NullDecimal {
	if !d.Valid {
		return Null
	}
	return d.Decimal.RoundBank(places).Nullable()
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
func (d NullDecimal) RoundCash(interval uint8) NullDecimal {
	if !d.Valid {
		return Null
	}
	return d.Decimal.RoundCash(interval).Nullable()
}

// Shift shifts the decimal in base 10.
// It shifts left when shift is positive and right if shift is negative.
// In simpler terms, the given value for shift is added to the exponent
// of the decimal.
func (d NullDecimal) Shift(shift int32) NullDecimal {
	if !d.Valid {
		return Null
	}
	return d.Decimal.Shift(shift).Nullable()
}

// Sin returns the sine of the radian argument x.
func (d NullDecimal) Sin() NullDecimal {
	if !d.Valid {
		return Null
	}
	return d.Decimal.Sin().Nullable()
}

// Sub returns d - d2.
func (d NullDecimal) Sub(d2 interface{}) NullDecimal {
	if v, ok := ValueOf(d2); ok {
		if d.Valid {
			return d.Decimal.Sub(v).Nullable()
		}
		return v.Neg().Nullable()
	}
	return d
}

// Tan returns the tangent of the radian argument x.
func (d NullDecimal) Tan() NullDecimal {
	if !d.Valid {
		return Null
	}
	return d.Decimal.Tan().Nullable()
}

// Truncate truncates off digits from the number, without rounding.
//
// NOTE: precision is the last digit that will not be truncated (must be >= 0).
//
// Example:
//
//     decimal.NewFromString("123.456").Truncate(2).String() // "123.45"
//
func (d NullDecimal) Truncate(precision int32) NullDecimal {
	if !d.Valid {
		return Null
	}
	return d.Decimal.Truncate(precision).Nullable()
}

//NullInt64 convert to sql.NullInt64
func (d NullDecimal) NullInt64() sql.NullInt64 {
	return sql.NullInt64{Int64: d.Decimal.IntPart(), Valid: d.Valid}
}

//NullFloat64 convert to sql.NullFloat64
func (d NullDecimal) NullFloat64() (f sql.NullFloat64) {
	f.Float64, _ = d.Decimal.Float64()
	f.Valid = d.Valid
	return
}

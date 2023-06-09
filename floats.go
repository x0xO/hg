package hg

import (
	"encoding/binary"
	"math"
	"math/big"
	"math/bits"
	"strconv"
)

// NewHFloat creates a new HFloat with the provided float64 value.
func NewHFloat(floats ...float64) HFloat {
	if len(floats) != 0 {
		return HFloat(floats[0])
	}

	return HFloat(0)
}

// Bytes returns the HFloat as a byte slice.
func (hf HFloat) Bytes() []byte {
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, hf.UInt64())

	return buffer[bits.LeadingZeros64(hf.UInt64())>>3:]
}

// Abs returns the absolute value of the HFloat.
func (hf HFloat) Abs() HFloat { return HFloat(math.Abs(hf.Float())) }

// Add adds two HFloats and returns the result.
func (hf HFloat) Add(b HFloat) HFloat { return hf + b }

// BigFloat returns the HFloat as a *big.Float.
func (hf HFloat) BigFloat() *big.Float { return big.NewFloat(hf.Float()) }

// Compare compares two HFloats and returns an HInt.
func (hf HFloat) Compare(b HFloat) HInt { return HInt(hf.BigFloat().Cmp(b.BigFloat())) }

// Div divides two HFloats and returns the result.
func (hf HFloat) Div(b HFloat) HFloat { return hf / b }

// Eq checks if two HFloats are equal.
func (hf HFloat) Eq(b HFloat) bool { return hf.Compare(b).Eq(0) }

// Float returns the HFloat as a float64.
func (hf HFloat) Float() float64 { return float64(hf) }

// Gt checks if the HFloat is greater than the specified HFloat.
func (hf HFloat) Gt(b HFloat) bool { return hf.Compare(b).Gt(0) }

// HInt returns the HFloat as an HInt.
func (hf HFloat) HInt() HInt { return HInt(hf) }

// HString returns the HFloat as an HString.
func (hf HFloat) HString() HString { return HString(strconv.FormatFloat(hf.Float(), 'f', -1, 64)) }

// Lt checks if the HFloat is less than the specified HFloat.
func (hf HFloat) Lt(b HFloat) bool { return hf.Compare(b).Lt(0) }

// Mul multiplies two HFloats and returns the result.
func (hf HFloat) Mul(b HFloat) HFloat { return hf * b }

// Ne checks if two HFloats are not equal.
func (hf HFloat) Ne(b HFloat) bool { return !hf.Eq(b) }

// Round rounds the HFloat to the nearest integer and returns the result as an HInt.
func (hf HFloat) Round() HInt { return HInt(math.Round(hf.Float())) }

// RoundDecimal rounds the HFloat value to the specified number of decimal places.
//
// The function takes the number of decimal places (precision) as an argument and returns a new
// HFloat value rounded to that number of decimals. This is achieved by multiplying the HFloat
// value by a power of 10 equal to the desired precision, rounding the result, and then dividing
// the rounded result by the same power of 10.
//
// Parameters:
//
// - precision (int): The number of decimal places to round the HFloat value to.
//
// Returns:
//
// - HFloat: A new HFloat value rounded to the specified number of decimal places.
//
// Example usage:
//
//	hf := hg.HFloat(3.14159)
//	rounded := hf.RoundDecimal(2) // rounded will be 3.14
func (hf HFloat) RoundDecimal(precision int) HFloat {
	mult := HFloat(math.Pow(10, float64(precision)))
	return hf.Mul(mult).Round().HFloat().Div(mult)
}

// Sub subtracts two HFloats and returns the result.
func (hf HFloat) Sub(b HFloat) HFloat { return hf - b }

// UInt64 returns the HFloat as a uint64.
func (hf HFloat) UInt64() uint64 { return math.Float64bits(hf.Float()) }

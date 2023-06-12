package hg

import (
	"encoding/binary"
	"fmt"
	"math/bits"
	"strconv"

	"github.com/x0xO/hg/pkg/rand"
)

// NewHInt creates a new HInt with the provided int value.
func NewHInt(ints ...int) HInt {
	if len(ints) != 0 {
		return HInt(ints[0])
	}

	return HInt(0)
}

// Min returns the minimum of two HInts.
func (hi HInt) Min(b HInt) HInt {
	if hi.Lt(b) {
		return hi
	}

	return b
}

// Max returns the maximum of two HInts.
func (hi HInt) Max(b HInt) HInt {
	if hi.Gt(b) {
		return hi
	}

	return b
}

// RandomRange returns a random HInt in the range [min, max].
func (HInt) RandomRange(min, max HInt) HInt {
	return HInt(rand.Intn(max.Sub(min).Add(1).Int())).Add(min)
}

// Bytes returns the HInt as a byte slice.
func (hi HInt) Bytes() []byte {
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, hi.UInt64())

	return buffer[bits.LeadingZeros64(hi.UInt64())>>3:]
}

// Add adds two HInts and returns the result.
func (hi HInt) Add(b HInt) HInt { return hi + b }

// Div divides two HInts and returns the result.
func (hi HInt) Div(b HInt) HInt { return hi / b }

// Eq checks if two HInts are equal.
func (hi HInt) Eq(b HInt) bool { return hi == b }

// Gt checks if the HInt is greater than the specified HInt.
func (hi HInt) Gt(b HInt) bool { return hi > b }

// Gte checks if the HInt is greater than or equal to the specified HInt.
func (hi HInt) Gte(b HInt) bool { return hi >= b }

// HFloat returns the HInt as an HFloat.
func (hi HInt) HFloat() HFloat { return HFloat(hi) }

// HString returns the HInt as an HString.
func (hi HInt) HString() HString { return HString(strconv.Itoa(int(hi))) }

// Int returns the HInt as an int.
func (hi HInt) Int() int { return int(hi) }

// Int16 returns the HInt as an int16.
func (hi HInt) Int16() int16 { return int16(hi) }

// Int32 returns the HInt as an int32.
func (hi HInt) Int32() int32 { return int32(hi) }

// Int64 returns the HInt as an int64.
func (hi HInt) Int64() int64 { return int64(hi) }

// Int8 returns the HInt as an int8.
func (hi HInt) Int8() int8 { return int8(hi) }

// IsNegative checks if the HInt is negative.
func (hi HInt) IsNegative() bool { return hi.Lt(0) }

// IsPositive checks if the HInt is positive.
func (hi HInt) IsPositive() bool { return hi.Gte(0) }

// Lt checks if the HInt is less than the specified HInt.
func (hi HInt) Lt(b HInt) bool { return hi < b }

// Lte checks if the HInt is less than or equal to the specified HInt.
func (hi HInt) Lte(b HInt) bool { return hi <= b }

// Mul multiplies two HInts and returns the result.
func (hi HInt) Mul(b HInt) HInt { return hi * b }

// Ne checks if two HInts are not equal.
func (hi HInt) Ne(b HInt) bool { return hi != b }

// Random returns a random HInt in the range [0, hi].
func (hi HInt) Random() HInt { return hi.RandomRange(0, hi) }

// Rem returns the remainder of the division between the receiver and the input value.
func (hi HInt) Rem(b HInt) HInt { return hi % b }

// Sub subtracts two HInts and returns the result.
func (hi HInt) Sub(b HInt) HInt { return hi - b }

// ToBinary returns the HInt as a binary string.
func (hi HInt) ToBinary() HString { return HString(fmt.Sprintf("%08b", hi)) }

// ToHex returns the HInt as a hexadecimal string.
func (hi HInt) ToHex() HString { return HString(fmt.Sprintf("%X", hi)) }

// ToOctal returns the HInt as an octal string.
func (hi HInt) ToOctal() HString { return HString(fmt.Sprintf("%o", hi)) }

// UInt returns the HInt as a uint.
func (hi HInt) UInt() uint { return uint(hi) }

// UInt16 returns the HInt as a uint16.
func (hi HInt) UInt16() uint16 { return uint16(hi) }

// UInt32 returns the HInt as a uint32.
func (hi HInt) UInt32() uint32 { return uint32(hi) }

// UInt64 returns the HInt as a uint64.
func (hi HInt) UInt64() uint64 { return uint64(hi) }

// UInt8 returns the HInt as a uint8.
func (hi HInt) UInt8() uint8 { return uint8(hi) }

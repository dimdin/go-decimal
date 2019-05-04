// Copyright 2014 Dimitris Dinodimos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package decimal

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// decomposer returns the internal decimal state into parts.
// If the provided// Decimal composes or decomposes a decimal value to and from individual parts.
// There are four separate parts: a boolean negative flag, a form byte with three possible states
// (finite=0, infinite=1, NaN=2),  a base-2 big-endian integer
// coefficient (also known as a significand) as a []byte, and an int32 exponent.
// These are composed into a final value as "decimal = (neg) (form=finite) coefficient * 10 ^ exponent".
// A zero length coefficient is a zero value.
// If the form is not finite the coefficient and scale should be ignored.
// The negative parameter may be set to true for any form, although implementations are not required
// to respect the negative parameter in the non-finite form.
//
// Implementations may choose to signal a negative zero or negative NaN, but implementations
// that do not support these may also ignore the negative zero or negative NaN without error.
// If an implementation does not support Infinity it may be converted into a NaN without error.
// If a value is set that is larger then what is supported by an implementation is attempted to
// be set, an error must be returned.
// Implementations must return an error if a NaN or Infinity is attempted to be set while neither
// are supported.
type decomposer interface {
	// Decompose returns the internal decimal state into parts.
	// If the provided buf has sufficient capacity, buf may be returned as the coefficient with
	// the value set and length set as appropriate.
	Decompose(buf []byte) (form byte, negative bool, coefficient []byte, exponent int32)

	// Compose sets the internal decimal value from parts. If the value cannot be
	// represented then an error should be returned.
	Compose(form byte, negative bool, coefficient []byte, exponent int32) error
}

// Decompose returns the internal decimal state into parts.
// If the provided buf has sufficient capacity, buf may be returned as the coefficient with
// the value set and length set as appropriate.
func (d Dec) Decompose(buf []byte) (form byte, negative bool, coefficient []byte, exponent int32) {
	high, low := d.coef.hi, d.coef.lo
	negative = high < 0
	exponent = -int32(d.scale)
	if negative {
		high = ^high
		low = -low
	}
	size := 16
	if high == 0 {
		size = 8
	}
	if cap(buf) >= size {
		coefficient = buf[:size]
	} else {
		coefficient = make([]byte, size)
	}
	if high == 0 {
		binary.BigEndian.PutUint64(coefficient, low)
	} else {
		binary.BigEndian.PutUint64(coefficient[8:], low)
		binary.BigEndian.PutUint64(coefficient, uint64(high))
	}
	return form, negative, coefficient, exponent
}

// Compose sets the internal decimal value from parts. If the value cannot be
// represented then an error should be returned.
func (d *Dec) Compose(form byte, negative bool, coefficient []byte, exponent int32) (err error) {
	if d == nil {
		return errors.New("Dec must not be nil")
	}
	if form != 0 {
		return errors.New("invalid form, form must be finite")
	}
	if exponent > 255 {
		return fmt.Errorf("exponent too large")
	}
	if exponent < -255 {
		return fmt.Errorf("exponent too small")
	}
	var low, high uint64
	maxi := len(coefficient) - 1
	for i := range coefficient {
		v := coefficient[maxi-i]
		if i < 8 {
			low |= uint64(v) << uint(i*8)
		} else if i < 16 {
			high |= uint64(v) << uint((i-8)*8)
		} else if v != 0 {
			return fmt.Errorf("coefficent too large")
		}
	}
	if high > math.MaxInt64 {
		return fmt.Errorf("coefficent too large")
	}
	i128 := Int128{
		hi: int64(high),
		lo: low,
	}
	if negative {
		i128.hi = ^i128.hi
		i128.lo = -i128.lo
	}
	d.coef = i128
	d.scale = 0

	if exponent > 0 {
		// Attempt to reduce coefficent. If not possible, return an error.
		d.coef.Mul(&(d.coef), exp10(uint8(exponent)))
	} else {
		d.scale = uint8(-exponent)
	}
	return nil
}

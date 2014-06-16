// Copyright 2014 Dimitris Dinodimos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package decimal

import "database/sql/driver"

// NullDec represents a decimal that may be null.
type NullDec struct {
	dec   Dec
	valid bool
}

// Null returns true if the value of d is null
func (d NullDec) Null() bool {
	return !d.valid
}

// SetNull sets null to d and returns d.
func (d *NullDec) SetNull() *NullDec {
	d.valid = false
	return d
}

// String returns the value of d.
// Returns an empty string if the value is null.
func (d NullDec) String() string {
	if d.Null() {
		return ""
	}
	return d.dec.String()
}

// Bytes returns the value of d.
// Returns nil if the value is null.
func (d NullDec) Bytes() []byte {
	if d.Null() {
		return nil
	}
	return d.dec.Bytes()
}

// Scan implements the database Scanner interface.
func (d *NullDec) Scan(value interface{}) error {
	if value == nil {
		d.valid = false
		return nil
	}
	err := d.dec.Scan(value)
	if err != nil {
		return err
	}
	d.valid = true
	return nil
}

// Value implements the database driver Valuer interface.
func (d NullDec) Value() (driver.Value, error) {
	if d.Null() {
		return nil, nil
	}
	return d.dec.Value()
}

// SetDec sets d to x and returns d.
func (d *NullDec) SetDec(x *Dec) *NullDec {
	d.valid = true
	d.dec.Set(x)
	return d
}

// Dec returns d.
func (d *NullDec) Dec() *Dec {
	if d.Null() {
		return nil
	}
	return &d.dec
}

// Set sets d to x and returns d.
func (d *NullDec) Set(x *NullDec) *NullDec {
	d.valid = x.valid
	d.dec.Set(&x.dec)
	return d
}

// SetFloat64 sets d to the value of f
func (d *NullDec) SetFloat64(f float64) error {
	err := d.dec.SetFloat64(f)
	if err != nil {
		return err
	}
	d.valid = true
	return nil
}

// SetInt64 sets d to the value of i and returns d.
func (d *NullDec) SetInt64(i int64) *NullDec {
	d.valid = true
	d.dec.SetInt64(i)
	return d
}

// SetInt128 sets d to the value of x and returns d.
func (d *NullDec) SetInt128(x *Int128) *NullDec {
	d.valid = true
	d.dec.SetInt128(x)
	return d
}

// SetString sets d to the value of s
func (d *NullDec) SetString(s string) error {
	if s == "" {
		d.SetNull()
		return nil
	}
	err := d.dec.SetString(s)
	if err != nil {
		return err
	}
	d.valid = true
	return nil
}

// SetBytes sets d to the value of buf
func (d *NullDec) SetBytes(buf []byte) error {
	if buf == nil {
		d.SetNull()
		return nil
	}
	err := d.dec.SetBytes(buf)
	if err != nil {
		return err
	}
	d.valid = true
	return nil
}

// Sign returns:
//
//      -1 if d <  0
//       0 if d == 0
//      +1 if d >  0
//
func (d NullDec) Sign() int {
	if d.Null() {
		return 0
	}
	return d.dec.Sign()
}

// Cmp compares x and y and returns:
//
//   -1 if x <  y
//    0 if x == y
//   +1 if x >  y
//
func (x NullDec) Cmp(y *NullDec) int {
	if x.Null() || y.Null() {
		return 0
	}
	return x.dec.Cmp(&y.dec)
}

// Abs sets d to |x| (the absolute value of x) and returns d.
func (d *NullDec) Abs(x *NullDec) *NullDec {
	if x.Null() {
		d.SetNull()
	} else {
		d.dec.Abs(&x.dec)
		d.valid = true
	}
	return d
}

// Neg sets d to -x and returns d.
func (d *NullDec) Neg(x *NullDec) *NullDec {
	if x.Null() {
		d.SetNull()
	} else {
		d.dec.Neg(&x.dec)
		d.valid = true
	}
	return d
}

// Add sets d to the sum x+y and returns d.
// The scale of d is the larger of the scales of the two operands.
func (d *NullDec) Add(x, y *NullDec) *NullDec {
	if x.Null() || y.Null() {
		d.SetNull()
	} else {
		d.dec.Add(&x.dec, &y.dec)
		d.valid = true
	}
	return d
}

// Sub sets d to the difference x-y and returns d.
// The scale of d is the larger of the scales of the two operands.
func (d *NullDec) Sub(x, y *NullDec) *NullDec {
	if x.Null() || y.Null() {
		d.SetNull()
	} else {
		d.dec.Sub(&x.dec, &y.dec)
		d.valid = true
	}
	return d
}

// Mul sets d to the product x*y and returns d.
// The scale of d is the sum of the scales of the two operands.
func (d *NullDec) Mul(x, y *NullDec) *NullDec {
	if x.Null() || y.Null() {
		d.SetNull()
	} else {
		d.dec.Mul(&x.dec, &y.dec)
		d.valid = true
	}
	return d
}

// Div sets d to the rounded quotient x/y and returns d.
// If y is zero panics with Division by zero.
// The resulting value is rounded half up to the given scale.
func (d *NullDec) Div(x, y *NullDec, scale uint8) *NullDec {
	if x.Null() || y.Null() {
		d.SetNull()
	} else {
		d.dec.Div(&x.dec, &y.dec, scale)
		d.valid = true
	}
	return d
}

// Round d half up to the given scale and returns d
func (d *NullDec) Round(scale uint8) *NullDec {
	if !d.Null() {
		d.dec.Round(scale)
	}
	return d
}

// Power sets d = x^n and returns d
func (d *NullDec) Power(x *NullDec, n int) *NullDec {
	if x.Null() {
		d.SetNull()
	} else {
		d.dec.Power(&x.dec, n)
		d.valid = true
	}
	return d
}

// Float64 returns the nearest float64 representation of d.
func (d NullDec) Float64() float64 {
	if d.Null() {
		return 0
	}
	return d.dec.Float64()
}

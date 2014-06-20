// Copyright 2014 Dimitris Dinodimos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package decimal implements the Dec and NullDec types suitable for financial and monetary calculations.
package decimal

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"io"
	"math"
	"strconv"
	"strings"
)

// Dec is represented as an 128 bit integer scaled by a power of ten.
type Dec struct {
	coef  Int128
	scale uint8
}

var decOne = &Dec{Int128{1, 0}, 0}

// Set sets d to x and returns d.
func (d *Dec) Set(x *Dec) *Dec {
	if d != x {
		d.coef = x.coef
		d.scale = x.scale
	}
	return d
}

// SetInt128 sets d to x and returns d.
func (d *Dec) SetInt128(x *Int128) *Dec {
	d.scale = 0
	d.coef = *x
	return d
}

// Sign returns:
//
//      -1 if d <  0
//       0 if d == 0
//      +1 if d >  0
//
func (d Dec) Sign() int {
	return d.coef.Sign()
}

func exp10(scale uint8) *Int128 {
	switch scale {
	case 0:
		return intOne
	case 1:
		return intTen
	case 2:
		return &Int128{100, 0}
	case 3:
		return &Int128{1000, 0}
	case 4:
		return &Int128{10000, 0}
	case 5:
		return &Int128{100000, 0}
	case 6:
		return &Int128{1000000, 0}
	case 7:
		return &Int128{10000000, 0}
	case 8:
		return &Int128{100000000, 0}
	}
	var z Int128
	z.Power(intTen, uint(scale))
	return &z
}

func (d *Dec) rescale(scale uint8) *Dec {
	if scale == d.scale {
		return d
	} else if scale < d.scale {
		panic("loss of precision")
	}
	z := d.coef
	z.Mul(&z, exp10(scale-d.scale))
	return &Dec{
		coef:  z,
		scale: scale,
	}
}

func maxscale(x, y *Dec) (*Dec, *Dec) {
	if x.scale == y.scale {
		return x, y
	}
	if x.scale > y.scale {
		return x, y.rescale(x.scale)
	}
	return x.rescale(y.scale), y
}

// Cmp compares x and y and returns:
//
//   -1 if x <  y
//    0 if x == y
//   +1 if x >  y
//
func (x Dec) Cmp(y *Dec) int {
	dx, dy := maxscale(&x, y)
	return dx.coef.Cmp(&dy.coef)
}

// Abs sets d to |x| (the absolute value of x) and returns d.
func (d *Dec) Abs(x *Dec) *Dec {
	d.coef.Abs(&x.coef)
	if d != x {
		d.scale = x.scale
	}
	return d
}

// Neg sets d to -x and returns d.
func (d *Dec) Neg(x *Dec) *Dec {
	d.coef.Neg(&x.coef)
	if d != x {
		d.scale = x.scale
	}
	return d
}

// Add sets d to the sum x+y and returns d.
// The scale of d is the larger of the scales of the two operands.
func (d *Dec) Add(x, y *Dec) *Dec {
	dx, dy := maxscale(x, y)
	if d != dx {
		d.scale = dx.scale
	}
	d.coef.Add(&dx.coef, &dy.coef)
	return d
}

// Sub sets d to the difference x-y and returns d.
// The scale of d is the larger of the scales of the two operands.
func (d *Dec) Sub(x, y *Dec) *Dec {
	dx, dy := maxscale(x, y)
	if d != dx {
		d.scale = dx.scale
	}
	d.coef.Sub(&dx.coef, &dy.coef)
	return d
}

// Mul sets d to the product x*y and returns d.
// The scale of d is the sum of the scales of the two operands.
func (d *Dec) Mul(x, y *Dec) *Dec {
	d.scale = x.scale + y.scale
	d.coef.Mul(&x.coef, &y.coef)
	return d
}

// Div sets d to the rounded quotient x/y and returns d.
// If y is zero panics with Division by zero.
// The resulting value is rounded half up to the given scale.
func (d *Dec) Div(x, y *Dec, scale uint8) *Dec {
	shift := int(scale) - int(x.scale) + int(y.scale)
	var sx, sy *Int128
	if shift > 0 {
		var z Int128
		sx = z.Mul(&x.coef, exp10(uint8(shift)))
		sy = &y.coef
	} else if shift < 0 {
		sx = &x.coef
		var z Int128
		sy = z.Mul(&y.coef, exp10(uint8(-shift)))
	} else {
		sx = &x.coef
		sy = &y.coef
	}
	d.scale = scale
	var r Int128
	d.coef.DivMod(sx, sy, &r)
	var roundUp bool
	if r.Sign() != 0 {
		r.Abs(&r)
		var v Int128
		v.Abs(sy)
		roundUp = r.Add(&r, &r).Cmp(&v) >= 0
	}
	if roundUp {
		if d.coef.Sign() < 0 {
			d.coef.Sub(&d.coef, intOne)
		} else {
			d.coef.Add(&d.coef, intOne)
		}
	}
	return d
}

// Round d half up to the given scale and returns d
func (d *Dec) Round(scale uint8) *Dec {
	if d.scale <= scale {
		return d
	}
	return d.Div(d, decOne, scale)
}

// Power sets d = x**n and returns d
func (d *Dec) Power(x *Dec, n int) *Dec {
	if n < 0 {
		scale := x.scale
		d.Power(x, -n)
		return d.Div(decOne, d, scale-uint8(n))
	} else if n == 0 {
		return d.Set(decOne)
	} else if n == 1 {
		return d.Set(x)
	} else if (n & 1) == 0 { // n even
		d.Mul(x, x)
		if d.scale > 18 {
			d.Round(18)
		}
		return d.Power(d, n/2)
	}
	// n odd
	var z Dec
	z.Set(x)
	d.Mul(x, x)
	if d.scale > 18 {
		d.Round(18)
	}
	d.Power(d, (n-1)/2)
	return d.Mul(d, &z)
}

// String returns the value of d
func (d Dec) String() string {
	return string(d.Bytes())
}

// Bytes returns the value of d
func (d Dec) Bytes() []byte {
	var dec Int128
	dec.Abs(&d.coef)
	digits := make([]byte, 0, 30)
	for i := 0; dec.Sign() != 0; i++ {
		var z Int128
		dec.DivMod(&dec, intTen, &z)
		digits = append(digits, byte(z.Int64()+'0'))
	}
	for int(d.scale) >= len(digits) {
		digits = append(digits, '0')
	}

	dst := make([]byte, 0, len(digits)+2)
	if d.Sign() < 0 {
		dst = append(dst, '-')
	}
	for i := len(digits); i != 0; i-- {
		if i == int(d.scale) {
			dst = append(dst, '.')
		}
		dst = append(dst, digits[i-1])
	}
	return dst
}

// SetFloat64 sets d to the value of f
func (d *Dec) SetFloat64(f float64) error {
	return d.SetString(strconv.FormatFloat(f, 'f', -1, 64))
}

// SetInt64 sets d to the value of i and returns d.
func (d *Dec) SetInt64(i int64) *Dec {
	d.coef.SetInt64(i)
	return d
}

// SetString sets d to the value of s
func (d *Dec) SetString(s string) error {
	if len(s) == 0 {
		return errors.New("SetString: empty string")
	}
	r := strings.NewReader(s)
	err := d.scan(r)
	if err != nil {
		return err
	}
	_, _, err = r.ReadRune()
	if err != io.EOF {
		return errors.New("SetString: non digit")
	}
	return nil
}

// SetBytes sets d to the value of buf
func (d *Dec) SetBytes(buf []byte) error {
	if len(buf) == 0 {
		return errors.New("SetBytes: empty buffer")
	}
	r := bytes.NewReader(buf)
	err := d.scan(r)
	if err != nil {
		return err
	}
	_, _, err = r.ReadRune()
	if err != io.EOF {
		return errors.New("SetBytes: non digit")
	}
	return nil
}

func (d *Dec) scan(r io.RuneScanner) error {
	d.coef.hi = 0
	d.coef.lo = 0
	d.scale = 0
	ch, _, err := r.ReadRune()
	if err != nil {
		return err
	}
	var neg bool
	switch ch {
	case '-':
		neg = true
	case '+':
	default:
		r.UnreadRune()
	}
	var dec bool
	for {
		ch, _, err = r.ReadRune()
		if err == io.EOF {
			goto ExitLoop
		}
		if err != nil {
			return err
		}
		switch {
		case ch == '.':
			if dec {
				r.UnreadRune()
				goto ExitLoop
			}
			dec = true
		case ch >= '0' && ch <= '9':
			d.coef.Mul(&d.coef, intTen)
			var z Int128
			z.SetInt64(int64(ch - '0'))
			d.coef.Add(&d.coef, &z)
			if dec {
				d.scale++
			}
		default:
			r.UnreadRune()
			goto ExitLoop
		}
	}
ExitLoop:
	if neg {
		d.Neg(d)
	}
	return nil
}

// Float64 returns the nearest float64 representation of d.
func (d Dec) Float64() float64 {
	return d.coef.Float64() / math.Pow10(int(d.scale))
}

// Scan implements the database Scanner interface.
func (d *Dec) Scan(value interface{}) error {
	if value == nil {
		return errors.New("Cannot Scan null into Dec")
	}
	switch value := value.(type) {
	case []byte:
		return d.SetBytes(value)
	case string:
		return d.SetString(value)
	case int64:
		d.SetInt64(value)
		return nil
	case float64:
		return d.SetFloat64(value)
	default:
		return errors.New("Invalid type Scan into Dec")
	}
}

// Value implements the database driver Valuer interface.
func (d Dec) Value() (driver.Value, error) {
	return d.Bytes(), nil
}

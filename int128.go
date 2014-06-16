// Copyright 2014 Dimitris Dinodimos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package decimal

// Int128 is a 128 bit signed integer.
type Int128 struct {
	lo uint64
	hi int64
}

var (
	intOne = &Int128{1, 0}
	intTen = &Int128{10, 0}
)

// SetInt64 sets z to x and returns z.
func (z *Int128) SetInt64(x int64) *Int128 {
	z.lo = uint64(x)
	if x < 0 {
		z.hi = -1
	} else {
		z.hi = 0
	}
	return z
}

// Set sets z to x and returns z.
func (z *Int128) Set(x *Int128) *Int128 {
	if z != x {
		z.lo = x.lo
		z.hi = x.hi
	}
	return z
}

// Int64 returns the int64 representation of x.
// If x cannot be represented in an int64, the result is undefined.
func (x Int128) Int64() int64 {
	return int64(x.lo)
}

// Float64 returns the nearest float64 representation of x.
func (x Int128) Float64() float64 {
	var d Int128
	d.Abs(&x)
	f := float64(d.lo) + float64(d.hi)*float64(1<<64)
	if x.Sign() < 0 {
		return -f
	}
	return f
}

// Sign returns:
//
//      -1 if x <  0
//       0 if x == 0
//      +1 if x >  0
//
func (x Int128) Sign() int {
	if x.lo == 0 && x.hi == 0 {
		return 0
	}
	if x.hi < 0 {
		return -1
	}
	return 1
}

// Cmp compares x and y and returns:
//
//   -1 if x <  y
//    0 if x == y
//   +1 if x >  y
//
func (x Int128) Cmp(y *Int128) int {
	if x.hi > y.hi {
		return 1
	} else if x.hi < y.hi {
		return -1
	} else if x.lo > y.lo {
		return 1
	} else if x.lo < y.lo {
		return -1
	}
	return 0
}

// Abs sets z to |x| (the absolute value of x) and returns z.
func (z *Int128) Abs(x *Int128) *Int128 {
	if x.hi < 0 {
		z.Neg(x)
	} else if z != x {
		z.hi = x.hi
		z.lo = x.lo
	}
	return z
}

// Neg sets z to -x and returns z.
func (z *Int128) Neg(x *Int128) *Int128 {
	if x.lo != 0 || x.hi != 0 {
		z.hi = ^x.hi
		z.lo = -x.lo
	} else if x != z {
		z.hi = 0
		z.lo = 0
	}
	return z
}

// Add sets z to the sum x+y and returns z.
func (z *Int128) Add(x, y *Int128) *Int128 {
	lo := x.lo
	z.lo = x.lo + y.lo
	z.hi = x.hi + y.hi
	if z.lo < lo {
		z.hi++
	}
	return z
}

// Sub sets z to the difference x-y and returns z.
func (z *Int128) Sub(x, y *Int128) *Int128 {
	lo := x.lo
	z.lo = x.lo - y.lo
	z.hi = x.hi - y.hi
	if z.lo > lo {
		z.hi--
	}
	return z
}

// Lsh sets z = x << n and returns z.
func (z *Int128) Lsh(x *Int128, n uint) *Int128 {
	if n > 63 {
		n -= 64
		z.hi = int64(x.lo)
		z.lo = 0
	} else {
		z.hi = x.hi
		z.lo = x.lo
	}
	if n == 0 {
		return z
	}
	z.hi = int64(uint64(z.hi)<<n + z.lo>>(64-n))
	z.lo <<= n
	return z
}

// Rsh sets z = x >> n and returns z.
func (z *Int128) Rsh(x *Int128, n uint) *Int128 {
	if n > 63 {
		n -= 64
		z.lo = uint64(z.hi)
		if x.hi < 0 {
			x.hi = -1
		} else {
			x.hi = 0
		}
	} else {
		z.hi = x.hi
		z.lo = x.lo
	}
	if n == 0 {
		return z
	}
	z.hi >>= n
	z.lo = uint64(z.hi)<<(64-n) + z.lo>>n
	return z
}

// Bit returns the value of the i'th bit of x. That is, it
// returns (x>>i)&1.
func (x *Int128) Bit(i int) uint {
	if i < 64 {
		return uint(x.lo>>uint(i)) & 1
	}
	return uint(x.hi>>uint(i-64)) & 1
}

// SetBit sets z to x, with x's i'th bit set to b (0 or 1).
// That is, if b is 1 SetBit sets z = x | (1 << i);
// if b is 0 SetBit sets z = x &^ (1 << i). If b is not 0 or 1,
// SetBit will panic.
func (z *Int128) SetBit(x *Int128, i int, b uint) *Int128 {
	if x != z {
		z.lo = x.lo
		z.hi = x.hi
	}
	if i < 64 {
		if b == 0 {
			z.lo = x.lo &^ (1 << uint(i))
		} else {
			z.lo = x.lo | (1 << uint(i))
		}
	} else {
		if b == 0 {
			z.hi = x.hi &^ (1 << uint(i-64))
		} else {
			z.hi = x.hi | (1 << uint(i-64))
		}
	}
	return z
}

const mask = 0xffffffff

// Algorithm M from Knuth TAOCP Vol 2 4.3.1
func mul(x, y, z *Int128) {
	var u, v [4]uint64
	var w [8]uint64
	u[0] = x.lo & mask
	u[1] = x.lo >> 32
	u[2] = uint64(x.hi) & mask
	u[3] = uint64(x.hi) >> 32

	v[0] = y.lo & mask
	v[1] = y.lo >> 32
	v[2] = uint64(y.hi) & mask
	v[3] = uint64(y.hi) >> 32

	var k, t uint64
	for j := 0; j < 4; j++ {
		// M3. Initialize i.
		k = 0
		for i := 0; i < 4; i++ {
			// M4. Multiply and add.
			t = u[i]*v[j] + w[i+j] + k
			// b = 2**32
			w[i+j] = t & mask // t mod b
			k = t >> 32       // t div b
		}
		w[j+4] = k
	}
	z.lo = w[0] | (w[1] << 32)
	z.hi = int64(w[2] | (w[3] << 32))
}

// Mul sets z to the product x*y and returns z.
func (z *Int128) Mul(x, y *Int128) *Int128 {
	if (x.lo == 0 && x.hi == 0) ||
		(y.lo == 0 && y.hi == 0) {
		z.lo = 0
		z.hi = 0
		return z
	}
	var u, v Int128
	u.Abs(x)
	v.Abs(y)
	mul(&u, &v, z)
	if (x.Sign() < 0) != (y.Sign() < 0) {
		z.Neg(z)
	}
	return z
}

func bits(x *Int128) (b uint) {
	var w uint64
	if x.hi != 0 {
		w = uint64(x.hi)
		b = 64
	} else {
		w = x.lo
	}
	for w != 0 {
		w >>= 1
		b++
	}
	return
}

// slow unsigned integer division
// TODO replace with Knuth algorithm D
func divmod(u, v, q, r *Int128) {
	q.lo = 0
	q.hi = 0
	r.lo = 0
	r.hi = 0
	n := bits(u)
	var i int
	for i = int(n) - 1; i >= 0; i-- {
		r.Lsh(r, 1)
		r.SetBit(r, 0, u.Bit(i))
		if r.Cmp(v) >= 0 {
			r.Sub(r, v)
			q.SetBit(q, i, 1)
		}
	}
}

func leadingZeros(x uint32) uint {
	if x == 0 {
		return 32
	}
	var n uint = 0
	if x <= 0x0000ffff {
		n += 16
		x <<= 16
	}
	if x <= 0x00ffffff {
		n += 8
		x <<= 8
	}
	if x <= 0x0fffffff {
		n += 4
		x <<= 4
	}
	if x <= 0x3fffffff {
		n += 2
		x <<= 2
	}
	if x <= 0x7fffffff {
		n++
	}
	return n
}

// Knuth TAOCP 4.3.1 algorithm D
func divmodD(xu, xv, xq, xr *Int128) {
	var v, q, r [4]uint64
	var u [5]uint64
	u[0] = xu.lo & mask
	u[1] = xu.lo >> 32
	u[2] = uint64(xu.hi) & mask
	u[3] = uint64(xu.hi) >> 32
	v[0] = xv.lo & mask
	v[1] = xv.lo >> 32
	v[2] = uint64(xv.hi) & mask
	v[3] = uint64(xv.hi) >> 32

	// D1. Normalize.
	n := 4
	for n >= 0 && v[n-1] == 0 {
		n--
	}
	shift := leadingZeros(uint32(v[n-1]))
	u[4] = u[3] >> (32 - shift)
	for i := 3; i > 0; i-- {
		u[i] = u[i]<<shift | u[i-1]>>(32-shift)
	}
	u[0] <<= shift
	for i := n - 1; i > 0; i-- {
		v[i] = v[i]<<shift | v[i-1]>>(32-shift)
	}
	v[0] <<= shift

	var b uint64 = 2 << 32
	var qhat, rhat, k, t, p uint64

	// D2. Initialize j. D7. Loop on j.
	for j := 4 - n; j >= 0; j-- {
		// D3. Calculate qhat.
		if u[j+n] == v[n-1] {
			qhat = b - 1
		} else {
			qhat = (u[j+n]*b + u[j+n-1]) / v[n-1]
		}
		rhat = (u[j+n]*b + u[j+n-1]) - qhat*v[n-1]
		for v[n-2]*qhat > rhat*b+u[j+n-2] {
			qhat--
			rhat += v[n-1]
			if rhat >= b {
				break
			}
		}

		// D4. Multiply and subtract.
		// u = u - qhat*v
		// M3. Initialize i.
		k = 0
		for i := 0; i < n; i++ {
			// M4. Multiply and subtract.
			p = qhat * v[i]
			t = u[i+j] - k - (p & mask)
			u[i+j] = t & mask
			k = (p >> 32) - t>>32
		}
		t = u[j+n] - k
		u[j+n] = t & mask

		// D5. Test remainder.
		q[j] = qhat
		if int64(t) < 0 {
			// D6. Add back.
			q[j]--
			// add v to u
			k = 0
			for i := 0; i < n; i++ {
				t = u[i+j] + v[i] + k
				u[i+j] = t & mask
				k = t >> 32
			}
			u[j+n] += k
		}
	}

	// D8. Unnormalize.
	for i := 0; i < n-1; i++ {
		r[i] = u[i]>>shift | u[i+1]<<(32-shift)
	}

	r[n-1] = u[n-1] >> shift
	xq.hi = int64(q[2] | (q[3] << 32))
	xq.lo = q[0] | (q[1] << 32)
	xr.hi = int64(r[2] | (r[3] << 32))
	xr.lo = r[0] | (r[1] << 32)
	return
}

// Knuth TAOCP 4.3.1 Exercise 16 algorithm
func divmod32(xu, xv, xq, xr *Int128) {
	var u, q [4]uint64
	u[0] = xu.lo & mask
	u[1] = xu.lo >> 32
	u[2] = uint64(xu.hi) & mask
	u[3] = uint64(xu.hi) >> 32

	var r uint64
	for j := 3; j >= 0; j-- {
		// r<<32 = r*b where b=2**32
		q[j] = (r<<32 + u[j]) / xv.lo
		r = (r<<32 + u[j]) % xv.lo
	}
	xq.hi = int64(q[2] | (q[3] << 32))
	xq.lo = q[0] | (q[1] << 32)
	xr.hi = 0
	xr.lo = r
	return
}

// DivMod sets z to the quotient x/y and r to the modulus x%y
// and returns the pair (z, r) for y != 0.
// If y == 0, a division-by-zero run-time panic occurs.
func (z *Int128) DivMod(x, y, r *Int128) (*Int128, *Int128) {
	if y.hi == 0 && y.lo == 0 {
		panic("Division by zero")
	} else if y.lo == 1 && y.hi == 0 {
		if z != x {
			z.lo = x.lo
			z.hi = x.hi
		}
		r.lo = 0
		r.hi = 0
		return z, r
	} else if x.lo == 0 && x.hi == 0 {
		z.lo = 0
		z.hi = 0
		r.lo = 0
		r.hi = 0
		return z, r
	} else if x == y || (x.lo == y.lo && x.hi == y.hi) {
		z.lo = 1
		z.hi = 0
		r.lo = 0
		r.hi = 0
		return z, r
	}
	var u, v, q Int128
	u.Abs(x)
	v.Abs(y)
	if v.Cmp(&u) > 0 {
		if r != x {
			r.lo = x.lo
			r.hi = x.hi
		}
		z.lo = 0
		z.hi = 0
		return z, r
	}

	if u.hi == 0 && v.hi == 0 {
		q.lo = u.lo / v.lo
		r.lo = u.lo % v.lo
	} else if v.hi == 0 && (v.lo>>32) == 0 {
		divmod32(&u, &v, &q, r)
	} else {
		divmod(&u, &v, &q, r)
	}

	if (x.Sign() < 0) != (y.Sign() < 0) {
		q.Neg(&q)
	}
	if x.Sign() < 0 {
		r.Neg(r)
	}
	z.lo = q.lo
	z.hi = q.hi
	return z, r
}

// Div sets z to the quotient x/y and returns z.
func (z *Int128) Div(x, y *Int128) *Int128 {
	var r Int128
	z.DivMod(x, y, &r)
	return z
}

// Mod sets z to the modulus x/y and returns z.
func (z *Int128) Mod(x, y *Int128) *Int128 {
	var q Int128
	q.DivMod(x, y, z)
	return z
}

// Power sets z = x^n and returns z
func (z *Int128) Power(x *Int128, n uint) *Int128 {
	if n == 0 {
		return z.Set(intOne)
	} else if n == 1 {
		return z.Set(x)
	} else if (n % 2) == 0 { // even
		z.Mul(x, x)
		return z.Power(z, n/2)
	}
	// odd
	var t Int128
	t.Set(x)
	z.Mul(x, x)
	z.Power(z, (n-1)/2)
	return z.Mul(z, &t)
}

// String returns the value of i
func (i Int128) String() string {
	return string(i.Bytes())
}

// Bytes returns the value of i
func (i Int128) Bytes() []byte {
	var dec Int128
	dec.Abs(&i)
	digits := make([]byte, 0, 30)
	for i := 0; dec.Sign() != 0; i++ {
		var z Int128
		dec.DivMod(&dec, intTen, &z)
		digits = append(digits, byte(z.Int64()+'0'))
	}
	dst := make([]byte, 0, len(digits)+1)
	if i.Sign() < 0 {
		dst = append(dst, '-')
	}
	for i := len(digits); i != 0; i-- {
		dst = append(dst, digits[i-1])
	}
	return dst
}

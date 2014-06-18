// Copyright 2014 Dimitris Dinodimos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package decimal

import (
	"math"
	"math/rand"
	"testing"
)

func TestIntSign(t *testing.T) {
	values := []struct {
		x int64
		s int
	}{
		{-1, -1},
		{-10, -1},
		{-100, -1},
		{0, 0},
		{1, 1},
		{10, 1},
		{100, 1},
	}
	for _, v := range values {
		var x Int128
		x.SetInt64(v.x)
		s := x.Sign()
		if v.s != s {
			t.Errorf("Sign %v does not match %v it is %v",
				v.x, v.s, s)
		}
	}
}

func TestIntNeg(t *testing.T) {
	values := []struct {
		s int64
		d int64
	}{
		{-1, 1},
		{-10, 10},
		{-100, 100},
		{0, 0},
		{1, -1},
		{10, -10},
		{100, -100},
	}
	for _, v := range values {
		var x, y Int128
		x.SetInt64(v.s)
		y.Neg(&x)
		if v.d != y.Int64() {
			t.Errorf("Neg %v does not match %v it is hi=%v lo=%v",
				v.s, v.d, y.hi, y.lo)
		}
	}
}

func TestIntAbs(t *testing.T) {
	values := []struct {
		s int64
		d int64
	}{
		{0, 0},
		{1, 1},
		{-1, 1},
		{10, 10},
		{-10, 10},
		{100, 100},
		{-100, 100},
	}
	for _, v := range values {
		var x, y Int128
		x.SetInt64(v.s)
		y.Abs(&x)
		if v.d != y.Int64() {
			t.Errorf("Abs %v does not match %v it is hi=%v lo=%v",
				v.s, v.d, y.hi, y.lo)
		}
	}
}

func TestIntCmp(t *testing.T) {
	values := []struct {
		x int64
		y int64
		c int
	}{
		{0, 0, 0},
		{1, 0, 1},
		{10, 0, 1},
		{-1, 0, -1},
		{-10, 0, -1},
		{0, 1, -1},
		{0, 10, -1},
		{0, -1, 1},
		{0, -10, 1},
		{1, 1, 0},
		{-1, 1, -1},
		{1, -1, 1},
		{10, 1, 1},
		{1, 10, -1},
	}
	for _, v := range values {
		var x, y Int128
		x.SetInt64(v.x)
		y.SetInt64(v.y)
		c := x.Cmp(&y)
		if v.c != c {
			t.Errorf("Cmp %v with %v does not match %v it is %v",
				v.x, v.y, v.c, c)
		}
	}
}

func TestIntAdd(t *testing.T) {
	values := []struct {
		x int64
		y int64
		z int64
	}{
		{0, 0, 0},
		{1, 0, 1},
		{10, 0, 10},
		{-1, 0, -1},
		{-10, 0, -10},
		{0, 1, 1},
		{0, 10, 10},
		{0, -1, -1},
		{0, -10, -10},
		{1, 1, 2},
		{-1, 1, 0},
		{1, -1, 0},
		{10, 1, 11},
		{1, 10, 11},
	}
	for _, v := range values {
		var x, y, z Int128
		x.SetInt64(v.x)
		y.SetInt64(v.y)
		z.Add(&x, &y)
		if v.z != z.Int64() {
			t.Errorf("Add %v with %v does not match %v it is %v",
				v.x, v.y, v.z, z.Int64())
		}
	}
}

func TestIntAddCarry(t *testing.T) {
	var x Int128
	x.lo = math.MaxUint64
	x.hi = 0
	var one Int128
	one.lo = 1
	one.hi = 0
	var z Int128
	z.Add(&x, &one)
	if z.lo != 0 && z.hi != 1 {
		t.Errorf("Got lo=%v hi=%v", z.lo, z.hi)
	}
}

func TestIntSub(t *testing.T) {
	values := []struct {
		x int64
		y int64
		z int64
	}{
		{0, 0, 0},
		{1, 0, 1},
		{10, 0, 10},
		{-1, 0, -1},
		{-10, 0, -10},
		{0, 1, -1},
		{0, 10, -10},
		{0, -1, 1},
		{0, -10, 10},
		{1, 1, 0},
		{-1, 1, -2},
		{1, -1, 2},
		{-1, -1, 0},
		{10, 1, 9},
		{1, 10, -9},
	}
	for _, v := range values {
		var x, y, z Int128
		x.SetInt64(v.x)
		y.SetInt64(v.y)
		z.Sub(&x, &y)
		if v.z != z.Int64() {
			t.Errorf("Sub from %v the %v does not match %v it is %v",
				v.x, v.y, v.z, z.Int64())
		}
	}
}

func TestIntSubBorrow(t *testing.T) {
	var x Int128
	x.lo = 0
	x.hi = 1
	var one Int128
	one.lo = 1
	one.hi = 0
	var z Int128
	z.Sub(&x, &one)
	if z.lo != math.MaxUint64 && z.hi != 0 {
		t.Errorf("Got lo=%v hi=%v", z.lo, z.hi)
	}
}

func TestLsh(t *testing.T) {
	values := []struct {
		x int64
		n uint
		z int64
	}{
		{0, 0, 0},
		{1, 0, 1},
		{1, 1, 2},
		{1, 10, 1024},
		{2, 2, 8},
	}
	for _, v := range values {
		var x, z Int128
		x.SetInt64(v.x)
		z.Lsh(&x, v.n)
		if v.z != z.Int64() {
			t.Errorf("Lsh %v << %v does not match %v it is %v",
				v.x, v.n, v.z, z.Int64())
		}
	}
}

func TestRsh(t *testing.T) {
	values := []struct {
		x int64
		n uint
		z int64
	}{
		{0, 0, 0},
		{1, 0, 1},
		{2, 1, 1},
		{1024, 10, 1},
		{8, 2, 2},
	}
	for _, v := range values {
		var x, z Int128
		x.SetInt64(v.x)
		z.Rsh(&x, v.n)
		if v.z != z.Int64() {
			t.Errorf("Rsh %v >> %v does not match %v it is %v",
				v.x, v.n, v.z, z.Int64())
		}
	}
}

func TestRandIntMul(t *testing.T) {
	var x, y, z, z2, acc Int128
	x.lo = uint64(rand.Int63())
	y.lo = uint64(rand.Int63() % 0xff)
	z.Mul(&x, &y)
	z2.Mul(&y, &x)
	var i uint64
	for i = 0; i < y.lo; i++ {
		acc.Add(&acc, &x)
	}
	if acc.Cmp(&z) != 0 {
		t.Errorf("Invalid multiplication %d x %d got %d,%d, want %d,%d", x.lo, y.lo, z.lo, z.hi, acc.lo, acc.hi)
	}
	if z2.Cmp(&z) != 0 {
		t.Errorf("Invalid multiplication")
	}
}

func TestIntMul(t *testing.T) {
	values := []struct {
		x int64
		y int64
		z int64
	}{
		{0, 0, 0},
		{2, 0, 0},
		{0, 2, 0},
		{1, 10, 10},
		{10, 1, 10},
		{4, 2, 8},
		{-4, 2, -8},
		{4, -2, -8},
		{-4, -2, 8},
	}
	for _, v := range values {
		var x, y, z Int128
		x.SetInt64(v.x)
		y.SetInt64(v.y)
		z.Mul(&x, &y)
		if v.z != z.Int64() {
			t.Errorf("Mul %v x %v does not match %v it is %v",
				v.x, v.y, v.z, z.Int64())
		}
		if x.Int64() != v.x {
			t.Errorf("Mul %v * %v alters %v to %v", v.x, v.y, v.x, x.Int64())
		}
		if y.Int64() != v.y {
			t.Errorf("Mul %v * %v alters %v to %v", v.x, v.y, v.y, y.Int64())
		}
	}
}

func TestRandIntDiv128(t *testing.T) {
	var x, y, q, r Int128
	x.lo = uint64(rand.Int63())
	x.hi = int64(rand.Int63())
	y.lo = uint64(rand.Int63())
	q.DivMod(&x, &y, &r)
	var tx Int128
	tx.Mul(&q, &y)
	tx.Add(&tx, &r)
	if tx.Cmp(&x) != 0 {
		t.Errorf("Invalid division %s / %s got q=%s r=%s", x, y, q, r)
	}
}

func TestRandIntDiv64(t *testing.T) {
	var x, y, q, r Int128
	x.lo = uint64(rand.Int63())
	y.lo = uint64(rand.Int63())
	q.DivMod(&x, &y, &r)
	var tx Int128
	tx.Mul(&q, &y)
	tx.Add(&tx, &r)
	if tx.Cmp(&x) != 0 {
		t.Errorf("Invalid division %d,%d / %d,%d got %d,%d rem %d,%d", x.lo, x.hi, y.lo, y.hi, q.lo, q.hi, r.lo, r.hi)
	}
}

func TestRandIntDiv10(t *testing.T) {
	var x, y, q, r Int128
	x.lo = uint64(rand.Int63())
	x.hi = int64(rand.Int63())
	y.lo = 10
	q.DivMod(&x, &y, &r)
	var tx Int128
	tx.Mul(&q, &y)
	tx.Add(&tx, &r)
	if tx.Cmp(&x) != 0 {
		t.Errorf("Invalid division %d,%d / %d,%d got %d,%d rem %d,%d", x.lo, x.hi, y.lo, y.hi, q.lo, q.hi, r.lo, r.hi)
	}
}

func TestIntDivMod(t *testing.T) {
	values := []struct {
		x int64
		y int64
		z int64
		r int64
	}{
		{0, 1, 0, 0},
		{2, 1, 2, 0},
		{3, 2, 1, 1},
		{1, 10, 0, 1},
		{10, 1, 10, 0},
		{8, 2, 4, 0},
		{-8, 2, -4, 0},
		{8, -2, -4, 0},
		{-8, -2, 4, 0},
		{9, 2, 4, 1},
		{10, 10, 1, 0},
		{100, 10, 10, 0},
		{1000, 10, 100, 0},
		{10000, 10, 1000, 0},
		{100000, 10, 10000, 0},
		{1000000, 10, 100000, 0},
		{10000000, 10, 1000000, 0},
		{100000000, 10, 10000000, 0},
		{1000000000, 10, 100000000, 0},
		{10000000000, 10, 1000000000, 0},
		{100000000000, 10, 10000000000, 0},
		{1000000000000, 10, 100000000000, 0},
		{10000000000000, 10, 1000000000000, 0},
		{100000000000000, 10, 10000000000000, 0},
		{1000000000000000, 10, 100000000000000, 0},
		{10000000000000000, 10, 1000000000000000, 0},
		{100000000000000000, 10, 10000000000000000, 0},
	}
	for _, v := range values {
		var x, y, z, r Int128
		x.SetInt64(v.x)
		y.SetInt64(v.y)
		z.DivMod(&x, &y, &r)
		if v.z != z.Int64() || v.r != r.Int64() {
			t.Errorf("DivMod %v/%v does not match %v and %v it is %v and %v",
				v.x, v.y, v.z, v.r, z.Int64(), r.Int64())
		}
		if x.Int64() != v.x {
			t.Errorf("DivMod %v/%v alters %v to %v", v.x, v.y, v.x, x.Int64())
		}
		if y.Int64() != v.y {
			t.Errorf("DivMod %v/%v alters %v to %v", v.x, v.y, v.y, y.Int64())
		}
	}
}

func TestDivision(t *testing.T) {
	values := []struct {
		ulo uint64
		uhi uint64
		vlo uint64
		vhi uint64
		qlo uint64
		qhi uint64
		rlo uint64
		rhi uint64
	}{
		{3, 0, 2, 0, 1, 0, 1, 0},                                                       //0
		{3, 0, 3, 0, 1, 0, 0, 0},                                                       //1
		{3, 0, 4, 0, 0, 0, 3, 0},                                                       //2
		{0, 0, 0xffffffff, 0, 0, 0, 0, 0},                                              //3
		{0xffffffff, 0, 1, 0, 0xffffffff, 0, 0, 0},                                     //4
		{0xffffffff, 0, 0xffffffff, 0, 1, 0, 0, 0},                                     //5
		{0xffffffff, 0, 3, 0, 0x55555555, 0, 0, 0},                                     //6
		{0xffffffff, 0xffffffff, 1, 0, 0xffffffff, 0xffffffff, 0, 0},                   //7
		{0xffffffff, 0xffffffff, 0xffffffff, 0, 1, 1, 0, 0},                            //8
		{0, 0, 0, 1, 0, 0, 0, 0},                                                       //9
		{0, 7, 0, 3, 2, 0, 0, 1},                                                       //10
		{5, 7, 0, 3, 2, 0, 5, 1},                                                       //11
		{0, 6, 0, 2, 3, 0, 0, 0},                                                       //12
		{0x80000000, 0, 0x40000001, 0, 0x00000001, 0, 0x3fffffff, 0},                   //13
		{0x0000789a, 0x0000bcde, 0x0000789a, 0x0000bcde, 1, 0, 0, 0},                   //14
		{0x0000789b, 0x0000bcde, 0x0000789a, 0x0000bcde, 1, 0, 1, 0},                   //15
		{0x00007899, 0x0000bcde, 0x0000789a, 0x0000bcde, 0, 0, 0x00007899, 0x0000bcde}, //16
		{0x0000ffff, 0x0000ffff, 0x0000ffff, 0x0000ffff, 1, 0, 0, 0},                   //17
		{0x0000ffff, 0x0000ffff, 0x00000000, 0x00000001, 0x0000ffff, 0, 0x0000ffff, 0}, //18
	}
	for i, a := range values {
		var u, v, q, r Int128
		u.lo = a.ulo
		u.hi = int64(a.uhi)
		v.lo = a.vlo
		v.hi = int64(a.vhi)
		q.DivMod(&u, &v, &r)
		if q.lo != a.qlo || q.hi != int64(a.qhi) ||
			r.lo != a.rlo || r.hi != int64(a.rhi) {
			t.Errorf("Failed %d %s / %s got q=%s r=%s", i,
				u, v, q, r)
		}
	}

}

func TestLeadingZeros(t *testing.T) {
	var x uint32 = (1 << 31)
	for i := 0; i <= 32; i++ {
		if int(leadingZeros(x)) != i {
			t.Errorf("Failed at %x: got %d want %d", x, leadingZeros(x), i)
		}
		x >>= 1
	}
}

func TestIntPower(t *testing.T) {
	values := []struct {
		x int64
		n uint
		r int64
	}{
		{2, 0, 1},
		{2, 1, 2},
		{2, 2, 4},
		{2, 3, 8},
		{2, 4, 16},
		{2, 5, 32},
		{2, 6, 64},
		{10, 0, 1},
		{10, 1, 10},
		{10, 2, 100},
		{10, 3, 1000},
		{10, 4, 10000},
		{10, 5, 100000},
		{10, 6, 1000000},
	}
	for _, a := range values {
		var x, y Int128
		x.SetInt64(a.x)
		y.Power(&x, a.n)
		if y.Int64() != a.r {
			t.Errorf("%d^%d got %d want %d", a.x, a.n, y.Int64(), a.r)
		}
	}
}

func panics(f func()) (b bool) {
	defer func() {
		if r := recover(); r != nil {
			b = true
		}
	}()
	f()
	return
}

func TestMulOverflow(t *testing.T) {
	if !panics(func() {
		var x Int128
		x.Set(intOne)
		for i := 0; i <= 38; i++ {
			x.Mul(&x, intTen)
		}
	}) {
		t.Errorf("failed to overflow mul")
	}
}

func TestDivZero(t *testing.T) {
	if !panics(func() {
		var x, y Int128
		x.Set(intOne)
		x.Div(&x, &y)
	}) {
		t.Errorf("failed to divide by zero")
	}
}

func TestAddOverflow(t *testing.T) {
	if !panics(func() {
		var x Int128
		x.lo = math.MaxUint64
		x.hi = math.MaxInt64
		x.Add(&x, intOne)
	}) {
		t.Errorf("failed to overflow add")
	}
}

func TestSubOverflow(t *testing.T) {
	if !panics(func() {
		var x Int128
		x.lo = 0
		x.hi = math.MinInt64
		x.Sub(&x, intOne)
	}) {
		t.Errorf("failed to overflow sub")
	}
}

// Copyright 2014 Dimitris Dinodimos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package decimal

import (
	"math/rand"
	"testing"
)

func BenchmarkAdd(b *testing.B) {
	var x, y Int128
	x.lo = uint64(rand.Int63())
	x.hi = rand.Int63()
	y.lo = uint64(rand.Int63())
	y.hi = rand.Int63()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var z Int128
		z.Add(&x, &y)
	}
}

func BenchmarkSub(b *testing.B) {
	var x, y Int128
	x.lo = uint64(rand.Int63())
	x.hi = rand.Int63()
	y.lo = uint64(rand.Int63())
	y.hi = rand.Int63()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var z Int128
		z.Sub(&x, &y)
	}
}

func BenchmarkMul(b *testing.B) {
	var x, y Int128
	x.lo = uint64(rand.Int63())
	x.hi = rand.Int63()
	y.lo = uint64(rand.Int63())
	y.hi = rand.Int63()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var z Int128
		z.Mul(&x, &y)
	}
}

func BenchmarkDiv10(b *testing.B) {
	var x, y Int128
	x.lo = uint64(rand.Int63())
	x.hi = rand.Int63()
	y.lo = 10
	y.hi = 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var z Int128
		z.Div(&x, &y)
	}
}

func BenchmarkDiv64(b *testing.B) {
	var x, y Int128
	x.lo = uint64(rand.Int63())
	x.hi = 0
	y.lo = uint64(rand.Int63())
	y.hi = 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var z Int128
		z.Div(&x, &y)
	}
}

func BenchmarkDiv128(b *testing.B) {
	var x, y Int128
	x.lo = uint64(rand.Int63())
	x.hi = rand.Int63()
	y.lo = uint64(rand.Int63())
	y.hi = 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var z Int128
		z.Div(&x, &y)
	}
}

func BenchmarkPower(b *testing.B) {
	var x Int128
	x.lo = uint64(rand.Int63())
	x.hi = 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Power(&x, 24)
	}
}

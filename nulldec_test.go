// Copyright 2014 Dimitris Dinodimos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package decimal

import (
	"strconv"
	"testing"
)

func TestSetBytes(t *testing.T) {
	values := [][]byte{
		nil,
		[]byte("0"),
		[]byte("10"),
		[]byte("-1"),
	}
	for _, a := range values {
		var x NullDec
		x.SetBytes(a)
		v := x.Bytes()
		if string(v) != string(a) {
			t.Errorf("Failed to SetBytes %s, got %s", string(a), string(v))
		}
	}
}

func TestNullSign(t *testing.T) {
	values := []struct {
		x string
		s int
	}{
		{"1.2", 1},
		{"-1.2", -1},
		{"0.0", 0},
		{"", 0},
	}
	for _, a := range values {
		var x NullDec
		x.SetString(a.x)
		if x.Sign() != a.s {
			t.Errorf("sign of %s got %d want %d", a.x, x.Sign(), a.s)
		}
	}
}

func TestNullNeg(t *testing.T) {
	values := []struct {
		x string
		y string
	}{
		{"", ""},
		{"0", "0"},
		{"-0.1", "0.1"},
		{"0.1", "-0.1"},
	}
	for _, a := range values {
		var x, y NullDec
		x.SetString(a.x)
		y.Neg(&x)
		if y.String() != a.y {
			t.Errorf("neg of %s got %s want %s", a.x, y.String(), a.y)
		}
	}
}

func TestNullAbs(t *testing.T) {
	values := []struct {
		x string
		y string
	}{
		{"1.2", "1.2"},
		{"-1.2", "1.2"},
		{"0", "0"},
		{"", ""},
	}
	for _, a := range values {
		var x, y NullDec
		x.SetString(a.x)
		y.Abs(&x)
		if y.String() != a.y {
			t.Errorf("abs of %s got %s want %s", a.x, y.String(), a.y)
		}
	}
}

func TestNullCmp(t *testing.T) {
	values := []struct {
		x string
		y string
		c int
	}{
		{"1.2", "1.2", 0},
		{"0.1", "0.01", 1},
		{"0.01", "0.1", -1},
		{"", "0.1", 0},
		{"0.1", "", 0},
	}
	for _, a := range values {
		var x, y NullDec
		x.SetString(a.x)
		y.SetString(a.y)
		c := x.Cmp(&y)
		if c != a.c {
			t.Errorf("cmp %s <> %s got %d want %d", a.x, a.y, c, a.c)
		}
	}
}

func TestNullAdd(t *testing.T) {
	values := []struct {
		x string
		y string
		z string
	}{
		{"1.2", "1.2", "2.4"},
		{"-1.2", "1.2", "0.0"},
		{"0", "0", "0"},
		{"", "1", ""},
		{"1", "", ""},
	}
	for _, a := range values {
		var x, y, z NullDec
		x.SetString(a.x)
		y.SetString(a.y)
		z.Add(&x, &y)
		if a.z != z.String() {
			t.Errorf("%s + %s got %s want %s", a.x, a.y, z.String(), a.z)
		}
	}
}

func TestNullSub(t *testing.T) {
	values := []struct {
		x string
		y string
		z string
	}{
		{"1.2", "1.2", "0.0"},
		{"-1.2", "1.2", "-2.4"},
		{"0", "0", "0"},
		{"", "1", ""},
		{"1", "", ""},
	}
	for _, a := range values {
		var x, y, z NullDec
		x.SetString(a.x)
		y.SetString(a.y)
		z.Sub(&x, &y)
		if a.z != z.String() {
			t.Errorf("%s - %s got %s want %s", a.x, a.y, z.String(), a.z)
		}
	}
}

func TestNullMul(t *testing.T) {
	values := []struct {
		x string
		y string
		z string
	}{
		{"1.2", "0", "0.0"},
		{"1.2", "1", "1.2"},
		{"1.2", "-1", "-1.2"},
		{"1.2", "2", "2.4"},
		{"1", "", ""},
		{"", "1", ""},
	}
	for _, a := range values {
		var x, y, z NullDec
		x.SetString(a.x)
		y.SetString(a.y)
		z.Mul(&x, &y)
		if a.z != z.String() {
			t.Errorf("%s * %s got %s want %s", a.x, a.y, z.String(), a.z)
		}
	}
}

func TestNullRound(t *testing.T) {
	values := []struct {
		x     string
		scale uint8
		r     string
	}{
		{"1.234", 2, "1.23"},
		{"", 2, ""},
	}
	for i, a := range values {
		var x NullDec
		x.SetString(a.x)
		x.Round(a.scale)
		if x.String() != a.r {
			t.Errorf("#%d got %s want %s", i, x, a.r)
		}
	}
}

func TestNullDiv(t *testing.T) {
	values := []struct {
		x, y  string
		scale uint8
		r     string
	}{
		{"10", "1", 2, "10.00"},
		{"10", "2", 2, "5.00"},
		{"", "2", 0, ""},
		{"10", "", 0, ""},
	}
	for _, a := range values {
		var x, y, z NullDec
		x.SetString(a.x)
		y.SetString(a.y)
		z.Div(&x, &y, a.scale)
		if z.String() != a.r {
			t.Errorf("%s / %s round %d got %s want %s", a.x, a.y, a.scale, z.String(), a.r)
		}
	}
}

func TestNullPower(t *testing.T) {
	values := []struct {
		x string
		n int
		r string
	}{
		{"2", 0, "1"},
		{"2", 1, "2"},
		{"2", 2, "4"},
		{"", 2, ""},
	}
	for _, a := range values {
		var x, y NullDec
		x.SetString(a.x)
		y.Power(&x, a.n)
		if y.String() != a.r {
			t.Errorf("%s^%d got %s want %s", a.x, a.n, y.String(), a.r)
		}
		if x.String() != a.x {
			t.Errorf("%s^%d alters %s to %s", a.x, a.n, a.x, x.String())
		}
	}
}

func TestNullFloat64(t *testing.T) {
	values := []string{
		"12.34",
		"-12.34",
		"0",
		"1",
		"-1",
	}
	for _, v := range values {
		var d NullDec
		d.SetString(v)
		s := strconv.FormatFloat(d.Float64(), 'f', -1, 64)
		if v != s {
			t.Errorf("expecting %s got %s", v, s)
		}
	}
}

func TestNullSetFloat64(t *testing.T) {
	values := []float64{
		// 17 significant digits
		9876543210987654000000000000,
		-9876543210987654000000000000,
		9876543210.987654,
		-9876543210.987654,
		1234567890123456700000000000,
		-1234567890123456700000000000,
	}
	for _, v := range values {
		var d NullDec
		d.SetFloat64(v)
		f := d.Float64()
		if v != f {
			t.Errorf("expecting %f got %f", v, f)
		}
	}
}

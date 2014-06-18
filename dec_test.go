// Copyright 2014 Dimitris Dinodimos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package decimal

import (
	"strconv"
	"testing"
)

func TestPi(t *testing.T) {
	// 38 digits of π
	pi := "3.1415926535897932384626433832795028842"
	var d Dec
	d.SetString(pi)
	if d.String() != pi {
		t.Errorf("Failed, expected %s got %s", pi, d.String())
	}
	if !panics(func() {
		d.SetString(pi + "0")
	}) {
		t.Errorf("Failed to overflow with 39 digits pi")
	}
}

func TestSetErrors(t *testing.T) {
	values := []string{
		"",
		"x",
		"-x",
		"+x",
		"0x",
		"1x",
		".x",
		".0x",
		"..",
	}
	var d Dec
	for _, v := range values {
		err := d.SetString(v)
		if err == nil {
			t.Errorf("Failed, expected error for SetString %s", v)
		}
		err = d.SetBytes([]byte(v))
		if err == nil {
			t.Errorf("Failed, expected error for SetBytes %s", v)
		}
	}
}

func TestExp10(t *testing.T) {
	var scale uint8
	var sd Int128
	sd.Set(intOne)
	for scale = 0; scale < 15; scale++ {
		d := exp10(scale)
		if sd.Cmp(d) != 0 {
			t.Errorf("Failed for scale %d, expected %s got %s",
				scale, sd, d)
		}
		sd.Mul(&sd, intTen)
	}
}

func TestSign(t *testing.T) {
	values := []struct {
		x string
		s int
	}{
		{"1.2", 1},
		{"-1.2", -1},
		{"0.0", 0},
		{"-0.0", 0},
		{"-0.1", -1},
		{"0.1", 1},
	}
	for _, a := range values {
		var x Dec
		x.SetString(a.x)
		if x.Sign() != a.s {
			t.Errorf("sign of %s got %d want %d", a.x, x.Sign(), a.s)
		}
	}
}

func TestNeg(t *testing.T) {
	values := []struct {
		x string
		y string
	}{
		{"1.2", "-1.2"},
		{"-1.2", "1.2"},
		{"0", "0"},
		{"-0", "0"},
		{"0.0", "0.0"},
		{"-0.0", "0.0"},
		{"0.00", "0.00"},
		{"-0.00", "0.00"},
		{"-0.1", "0.1"},
		{"0.1", "-0.1"},
		{"-0.01", "0.01"},
		{"0.01", "-0.01"},
		{"-0.10", "0.10"},
		{"0.10", "-0.10"},
	}
	for _, a := range values {
		var x, y Dec
		x.SetString(a.x)
		y.Neg(&x)
		if y.String() != a.y {
			t.Errorf("neg of %s got %s want %s", a.x, y.String(), a.y)
		}
	}
}

func TestAbs(t *testing.T) {
	values := []struct {
		x string
		y string
	}{
		{"1.2", "1.2"},
		{"-1.2", "1.2"},
		{"0", "0"},
		{"-0", "0"},
		{"0.0", "0.0"},
		{"-0.0", "0.0"},
		{"0.00", "0.00"},
		{"-0.00", "0.00"},
		{"-0.1", "0.1"},
		{"0.1", "0.1"},
		{"-0.01", "0.01"},
		{"0.01", "0.01"},
		{"-0.10", "0.10"},
		{"0.10", "0.10"},
	}
	for _, a := range values {
		var x, y Dec
		x.SetString(a.x)
		y.Abs(&x)
		if y.String() != a.y {
			t.Errorf("abs of %s got %s want %s", a.x, y.String(), a.y)
		}
	}
}

func TestCmp(t *testing.T) {
	values := []struct {
		x string
		y string
		c int
	}{
		{"1.2", "1.2", 0},
		{"-1.2", "1.2", -1},
		{"0", "0", 0},
		{"-0.00", "0.0", 0},
		{"-0.01", "0.01", -1},
		{"0.01", "0.01", 0},
		{"0.10", "-0.10", 1},
		{"0.1", "0.01", 1},
		{"0.01", "0.1", -1},
	}
	for _, a := range values {
		var x, y Dec
		x.SetString(a.x)
		y.SetString(a.y)
		c := x.Cmp(&y)
		if c != a.c {
			t.Errorf("cmp %s <> %s got %d want %d", a.x, a.y, c, a.c)
		}
	}
}

func TestAdd(t *testing.T) {
	values := []struct {
		x string
		y string
		z string
	}{
		{"1.2", "1.2", "2.4"},
		{"-1.2", "1.2", "0.0"},
		{"0", "0", "0"},
		{"-0.00", "0.0", "0.00"},
		{"-0.01", "0.01", "0.00"},
		{"0.01", "0.01", "0.02"},
		{"0.1", "0.01", "0.11"},
		{"-0.01", "0.1", "0.09"},
	}
	for _, a := range values {
		var x, y, z Dec
		x.SetString(a.x)
		y.SetString(a.y)
		z.Add(&x, &y)
		if a.z != z.String() {
			t.Errorf("%s + %s got %s want %s", a.x, a.y, z.String(), a.z)
		}
	}
}

func TestSub(t *testing.T) {
	values := []struct {
		x string
		y string
		z string
	}{
		{"1.2", "1.2", "0.0"},
		{"-1.2", "1.2", "-2.4"},
		{"0", "0", "0"},
		{"-0.00", "0.0", "0.00"},
		{"-0.01", "0.01", "-0.02"},
		{"0.01", "0.01", "0.00"},
		{"0.1", "0.01", "0.09"},
		{"-0.01", "0.1", "-0.11"},
	}
	for _, a := range values {
		var x, y, z Dec
		x.SetString(a.x)
		y.SetString(a.y)
		z.Sub(&x, &y)
		if a.z != z.String() {
			t.Errorf("%s - %s got %s want %s", a.x, a.y, z.String(), a.z)
		}
	}
}

func TestMul(t *testing.T) {
	values := []struct {
		x string
		y string
		z string
	}{
		{"1.2", "0", "0.0"},
		{"1.2", "1", "1.2"},
		{"1.2", "-1", "-1.2"},
		{"1.2", "2", "2.4"},
		{"1.2", "-2", "-2.4"},
		{"1.2", "10", "12.0"},
		{"-1.2", "0", "0.0"},
		{"-1.2", "1", "-1.2"},
		{"-1.2", "-1", "1.2"},
		{"-1.2", "2", "-2.4"},
		{"-1.2", "-2", "2.4"},
		{"-1.2", "10", "-12.0"},
	}
	for _, a := range values {
		var x, y, z Dec
		x.SetString(a.x)
		y.SetString(a.y)
		z.Mul(&x, &y)
		if a.z != z.String() {
			t.Errorf("%s * %s got %s want %s", a.x, a.y, z.String(), a.z)
		}
		if x.String() != a.x {
			t.Errorf("%s * %s alters %s to %s", a.x, a.y, a.x, x.String())
		}
		if y.String() != a.y {
			t.Errorf("%s * %s alters %s to %s", a.x, a.y, a.y, y.String())
		}
	}
}

func TestRound(t *testing.T) {
	values := []struct {
		x     string
		scale uint8
		r     string
	}{
		{"1.234", 2, "1.23"},
		{"1.235", 2, "1.24"},
		{"-1.234", 2, "-1.23"},
		{"-1.235", 2, "-1.24"},
		{"1.23", 2, "1.23"},
	}
	for i, a := range values {
		var x Dec
		x.SetString(a.x)
		x.Round(a.scale)
		if x.String() != a.r {
			t.Errorf("#%d got %s want %s", i, x, a.r)
		}
	}
}

func TestDiv(t *testing.T) {
	values := []struct {
		x, y  string
		scale uint8
		r     string
	}{
		{"10", "1", 2, "10.00"},
		{"10", "2", 2, "5.00"},
		{"-10", "1", 2, "-10.00"},
		{"-10", "2", 2, "-5.00"},
		{"10", "-1", 2, "-10.00"},
		{"10", "-2", 2, "-5.00"},
		{"-10", "-1", 2, "10.00"},
		{"-10", "-2", 2, "5.00"},
		{"1", "1", 2, "1.00"},
		{"1", "2", 2, "0.50"},
		{"1", "3", 2, "0.33"},
		{"1", "4", 2, "0.25"},
		{"1", "5", 2, "0.20"},
		{"1", "10", 1, "0.1"},
		{"1", "100", 2, "0.01"},
		{"1", "1000", 3, "0.001"},
		{"1", "10000", 4, "0.0001"},
		{"10", "2", 0, "5"},
	}
	for _, a := range values {
		var x, y, z Dec
		x.SetString(a.x)
		y.SetString(a.y)
		z.Div(&x, &y, a.scale)
		if z.String() != a.r {
			t.Errorf("%s / %s round %d got %s want %s", a.x, a.y, a.scale, z.String(), a.r)
		}
	}
}

func TestPower(t *testing.T) {
	values := []struct {
		x string
		n int
		r string
	}{
		{"2", 0, "1"},
		{"2", 1, "2"},
		{"2", 2, "4"},
		{"2", 3, "8"},
		{"2", 4, "16"},
		{"2", 5, "32"},
		{"2", 6, "64"},
		{"2", 7, "128"},
		{"2", 8, "256"},
		{"10", 0, "1"},
		{"10", 1, "10"},
		{"10", 2, "100"},
		{"10", 3, "1000"},
		{"10", 4, "10000"},
		{"10", 5, "100000"},
		{"10", 6, "1000000"},
		{"10", 7, "10000000"},
		{"10", 8, "100000000"},
		{"10", 9, "1000000000"},
		{"10", -1, "0.1"},
		{"10", -2, "0.01"},
		{"10", -3, "0.001"},
		{"10", -4, "0.0001"},
		{"10", -5, "0.00001"},
		{"10", -6, "0.000001"},
		{"10", -7, "0.0000001"},
		{"10", -8, "0.00000001"},
		{"10", -9, "0.000000001"},
	}
	for _, a := range values {
		var x, y Dec
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

func TestFloat64(t *testing.T) {
	values := []string{
		"12.34",
		"-12.34",
		"0",
		"1",
		"-1",
		// 17 significant digits
		"9876543210987654000000000000",
		"-9876543210987654000000000000",
		"9876543210.987654",
		"-9876543210.987654",
		"1234567890123456700000000000",
		"-1234567890123456700000000000",
	}
	for _, v := range values {
		var d Dec
		d.SetString(v)
		s := strconv.FormatFloat(d.Float64(), 'f', -1, 64)
		if v != s {
			t.Errorf("expecting %s got %s", v, s)
		}
	}
}

func TestSetFloat64(t *testing.T) {
	values := []float64{
		12.34,
		-12.34,
		0,
		1,
		-1,
		// 17 significant digits
		9876543210987654000000000000,
		-9876543210987654000000000000,
		9876543210.987654,
		-9876543210.987654,
		1234567890123456700000000000,
		-1234567890123456700000000000,
	}
	for _, v := range values {
		var d Dec
		d.SetFloat64(v)
		f := d.Float64()
		if v != f {
			t.Errorf("expecting %f got %f", v, f)
		}
	}
}

// Copyright 2014 Dimitris Dinodimos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package decimal

import (
	"fmt"
	"strconv"
	"testing"
)

func ExampleNullDec_SetString() {
	var d NullDec
	d.SetString("12.34")
	fmt.Println(d)
	// Output:
	// 12.34
}

func ExampleNullDec_SetBytes() {
	var d NullDec
	bytes := []byte("-12.34")
	d.SetBytes(bytes)
	fmt.Println(d)
	d.SetString("")
	if d.Null() {
		fmt.Println("null")
	}
	// Output:
	// -12.34
	// null
}

func ExampleNullDec_Null() {
	var d NullDec
	d.SetString("-12.34")
	if !d.Null() {
		fmt.Println(d)
	}
	// Output:
	// -12.34
}

func ExampleNullDec_SetInt128() {
	var i Int128
	i.SetInt64(100000)
	var d NullDec
	d.SetInt128(&i)
	fmt.Println(d)
	// Output:
	// 100000
}

func ExampleNullDec_SetDec() {
	var d Dec
	d.SetString("1.2")
	var nd NullDec
	nd.SetDec(&d)
	fmt.Println(nd)
	// Output:
	// 1.2
}

func ExampleNullDec_Add() {
	var x, y NullDec
	x.SetString("0.1")
	y.SetInt64(1)
	x.Add(&x, &y)
	fmt.Println(x)
	// Output:
	// 1.1
}

func ExampleNullDec_Mul() {
	var x, y NullDec
	x.SetString("1.1")
	y.SetInt64(2)
	x.Mul(&x, &y)
	fmt.Println(x)
	// Output:
	// 2.2
}

func ExampleNullDec_Div() {
	var x, y NullDec
	x.SetInt64(100)
	y.SetInt64(3)
	x.Div(&x, &y, 2)
	fmt.Println(x)
	// Output:
	// 33.33
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

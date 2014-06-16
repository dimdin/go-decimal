// Copyright 2014 Dimitris Dinodimos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package decimal

import "fmt"

func ExampleDec_Cmp() {
	var x, y Dec
	x.SetString("1.1")
	y.SetString("2.0")
	if x.Cmp(&y) < 0 {
		fmt.Printf("%s < %s", x, y)
	}
	// Output:
	// 1.1 < 2.0
}

func ExampleNullDec_Cmp() {
	var x, y NullDec
	x.SetString("1.1")
	y.SetString("2.0")
	if x.Cmp(&y) < 0 {
		fmt.Printf("%s < %s", x, y)
	}
	// Output:
	// 1.1 < 2.0
}

func ExampleDec_Neg() {
	var d Dec
	d.SetString("12.34")
	d.Neg(&d)
	fmt.Println(d)
	// Output:
	// -12.34

}

func ExampleNullDec_Neg() {
	var d NullDec
	d.SetString("12.34")
	d.Neg(&d)
	fmt.Println(d)
	// Output:
	// -12.34

}

func ExampleDec_SetString() {
	var d Dec
	d.SetString("-12.34")
	fmt.Println(d)
	// Output:
	// -12.34
}

func ExampleNullDec_SetString() {
	var d NullDec
	d.SetString("12.34")
	fmt.Println(d)
	// Output:
	// 12.34
}

func ExampleDec_SetBytes() {
	var d Dec
	bytes := []byte("+12.34")
	d.SetBytes(bytes)
	fmt.Println(d)
	// Output:
	// 12.34
}

func ExampleNullDec_SetBytes() {
	var d NullDec
	bytes := []byte("-12.34")
	d.SetBytes(bytes)
	fmt.Println(d)
	// Output:
	// -12.34
}

func ExampleNullDec_Null() {
	var d NullDec
	d.SetString("-12.34")
	if !d.Null() {
		fmt.Println(d)
	}
	d.SetString("")
	if d.Null() {
		fmt.Println("null")
	}
	// Output:
	// -12.34
	// null
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

func ExampleDec_SetInt128() {
	var i Int128
	i.SetInt64(100000)
	var d Dec
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

func ExampleNullDec_Dec() {
	var d NullDec
	d.SetString("1.2")
	fmt.Println(d.Dec())
	// Output:
	// 1.2
}

func ExampleDec_Add() {
	var x, y Dec
	x.SetString("0.1")
	y.SetInt64(1)
	x.Add(&x, &y)
	fmt.Println(x)
	// Output:
	// 1.1
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

func ExampleDec_Mul() {
	var x, y Dec
	x.SetString("1.1")
	y.SetInt64(2)
	x.Mul(&x, &y)
	fmt.Println(x)
	// Output:
	// 2.2
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

func ExampleDec_Div() {
	var x, y Dec
	x.SetInt64(100)
	y.SetInt64(3)
	x.Div(&x, &y, 2)
	fmt.Println(x)
	// Output:
	// 33.33
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

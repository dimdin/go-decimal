// Copyright 2014 Dimitris Dinodimos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package decimal

import (
	"testing"
)

func TestDecomposerRoundTrip(t *testing.T) {
	list := []struct {
		N string
		S string
		E bool
	}{
		{N: "Zero", S: "0"},
		{N: "Normal-1", S: "123.456"},
		{N: "Normal-2", S: "432000"},
		{N: "Normal-3", S: "0.00456"},
		{N: "Neg-1", S: "-123.456"},
		{N: "Big-1", S: "42535295865117307919086767873688.862721"},
		{N: "Big-2", S: "-42535295865117307919086767873688.862721"},
		{N: "ParseError", S: "NaN", E: true},
	}
	for _, item := range list {
		d := &Dec{}
		err := d.SetString(item.S)
		if err == nil && item.E {
			t.Fatal("expected error, got <nil>")
		}
		if err != nil && !item.E {
			t.Fatalf("did not expect error, got %v", err)
		}
		if item.E {
			return
		}
		set := &Dec{}
		err = set.Compose(d.Decompose(nil))
		if err == nil && item.E {
			t.Fatal("expected error, got <nil>")
		}
		if err != nil && !item.E {
			t.Fatalf("did not expect error, got %v", err)
		}
		if d.Cmp(set) != 0 {
			t.Fatalf("for %q, wanted %v, got %v", item.S, d, set)
		}
		if s := set.String(); s != item.S {
			t.Fatalf("wanted %q got %q", item.S, s)
		}
	}
}

func TestDecomposerCompose(t *testing.T) {
	list := []struct {
		N string // Name.
		S string // String value.

		Form byte // Form
		Neg  bool
		Coef []byte // Coefficent
		Exp  int32

		Err bool // Expect an error.
	}{
		{N: "Zero", S: "0", Coef: nil, Exp: 0},
		{N: "Normal-1", S: "123.456", Coef: []byte{0x01, 0xE2, 0x40}, Exp: -3},
		{N: "Neg-1", S: "-123.456", Neg: true, Coef: []byte{0x01, 0xE2, 0x40}, Exp: -3},
		{N: "PosExp-1", S: "123456000", Coef: []byte{0x01, 0xE2, 0x40}, Exp: 3},
		{N: "PosExp-2", S: "-123456000", Neg: true, Coef: []byte{0x01, 0xE2, 0x40}, Exp: 3},
		{N: "AllDec-1", S: "0.123456", Coef: []byte{0x01, 0xE2, 0x40}, Exp: -6},
		{N: "AllDec-2", S: "-0.123456", Neg: true, Coef: []byte{0x01, 0xE2, 0x40}, Exp: -6},
		{N: "Big-1", S: "42535295865117307919086767873688.862721", Neg: false, Coef: []byte{0x1F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, Exp: -6},
		{N: "Big-2", S: "-42535295865117307919086767873688.862721", Neg: true, Coef: []byte{0x1F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, Exp: -6},
	}

	for _, item := range list {
		d := &Dec{}
		err := d.SetString(item.S)
		if err != nil {
			t.Fatal(err)
		}
		err = d.Compose(item.Form, item.Neg, item.Coef, item.Exp)
		if err != nil && !item.Err {
			t.Fatalf("unexpected error, got %v", err)
		}
		if item.Err {
			if err == nil {
				t.Fatal("expected error, got <nil>")
			}
			return
		}
		if s := d.String(); s != item.S {
			t.Fatalf("unexpected value, got %q want %q", s, item.S)
		}
	}
}

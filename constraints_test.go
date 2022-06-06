//
// Copyright 2018-2022 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConstraints(t *testing.T) {
	lt := &LessThan{v("1.3.0")}
	require.True(t, lt.Match(v("1.0.0")))
	require.False(t, lt.Match(v("1.3.0")))
	require.False(t, lt.Match(v("2.0.0")))
	require.Equal(t, "<1.3.0", lt.String())

	lte := &LessThanOrEqual{v("1.3.0")}
	require.True(t, lte.Match(v("1.0.0")))
	require.True(t, lte.Match(v("1.3.0")))
	require.False(t, lte.Match(v("2.0.0")))
	require.Equal(t, "<=1.3.0", lte.String())

	eq := &Equals{v("1.3.0")}
	require.False(t, eq.Match(v("1.0.0")))
	require.True(t, eq.Match(v("1.3.0")))
	require.False(t, eq.Match(v("2.0.0")))
	require.Equal(t, "=1.3.0", eq.String())

	gte := &GreaterThanOrEqual{v("1.3.0")}
	require.False(t, gte.Match(v("1.0.0")))
	require.True(t, gte.Match(v("1.3.0")))
	require.True(t, gte.Match(v("2.0.0")))
	require.Equal(t, ">=1.3.0", gte.String())

	gt := &GreaterThan{v("1.3.0")}
	require.False(t, gt.Match(v("1.0.0")))
	require.False(t, gt.Match(v("1.3.0")))
	require.True(t, gt.Match(v("2.0.0")))
	require.Equal(t, ">1.3.0", gt.String())

	tr := &True{}
	require.True(t, tr.Match(v("1.0.0")))
	require.True(t, tr.Match(v("1.3.0")))
	require.True(t, tr.Match(v("2.0.0")))
	require.Equal(t, "", tr.String())

	gt100 := &GreaterThan{v("1.0.0")}
	lte200 := &LessThanOrEqual{v("2.0.0")}
	and := &And{[]Constraint{gt100, lte200}}
	require.False(t, and.Match(v("0.9.0")))
	require.False(t, and.Match(v("1.0.0")))
	require.True(t, and.Match(v("1.3.0")))
	require.True(t, and.Match(v("2.0.0")))
	require.False(t, and.Match(v("2.1.0")))
	require.Equal(t, "(>1.0.0 && <=2.0.0)", and.String())

	gt200 := &GreaterThan{v("2.0.0")}
	lte100 := &LessThanOrEqual{v("1.0.0")}
	or := &Or{[]Constraint{gt200, lte100}}
	require.True(t, or.Match(v("0.9.0")))
	require.True(t, or.Match(v("1.0.0")))
	require.False(t, or.Match(v("1.3.0")))
	require.False(t, or.Match(v("2.0.0")))
	require.True(t, or.Match(v("2.1.0")))
	require.Equal(t, "(>2.0.0 || <=1.0.0)", or.String())

	notOr := &Not{or}
	require.False(t, notOr.Match(v("0.9.0")))
	require.False(t, notOr.Match(v("1.0.0")))
	require.True(t, notOr.Match(v("1.3.0")))
	require.True(t, notOr.Match(v("2.0.0")))
	require.False(t, notOr.Match(v("2.1.0")))
	require.Equal(t, "!(>2.0.0 || <=1.0.0)", notOr.String())

	comp := &CompatibleWith{v("1.3.4-rc.3")}
	require.False(t, comp.Match(v("1.2.3")))
	require.False(t, comp.Match(v("1.3.2")))
	require.False(t, comp.Match(v("1.2.3")))
	require.False(t, comp.Match(v("1.3.4-rc.1")))
	require.True(t, comp.Match(v("1.3.4-rc.5")))
	require.True(t, comp.Match(v("1.3.4")))
	require.True(t, comp.Match(v("1.3.6")))
	require.True(t, comp.Match(v("1.4.0")))
	require.True(t, comp.Match(v("1.4.5")))
	require.True(t, comp.Match(v("1.4.5-rc.2")))
	require.False(t, comp.Match(v("2.0.0")))
}

func TestConstraintsParser(t *testing.T) {
	type goodStringTest struct {
		In, Out string
	}
	good := []goodStringTest{
		{"", ""}, // always true
		{"=1.3.0", "=1.3.0"},
		{" =1.3.0 ", "=1.3.0"},
		{"=1.3.0 ", "=1.3.0"},
		{" =1.3.0", "=1.3.0"},
		{">=1.3.0", ">=1.3.0"},
		{">1.3.0", ">1.3.0"},
		{"<=1.3.0", "<=1.3.0"},
		{"<1.3.0", "<1.3.0"},
		{"^1.3.0", "^1.3.0"},
		{" ^1.3.0", "^1.3.0"},
		{"^1.3.0 ", "^1.3.0"},
		{" ^1.3.0 ", "^1.3.0"},
		{"(=1.4.0)", "=1.4.0"},
		{"!(=1.4.0)", "!(=1.4.0)"},
		{"!(((=1.4.0)))", "!(=1.4.0)"},
		{"=1.2.4 && =1.3.0", "(=1.2.4 && =1.3.0)"},
		{"=1.2.4 && ^1.3.0", "(=1.2.4 && ^1.3.0)"},
		{"=1.2.4 && =1.3.0 && =1.2.0", "(=1.2.4 && =1.3.0 && =1.2.0)"},
		{"=1.2.4 && =1.3.0 || =1.2.0", "((=1.2.4 && =1.3.0) || =1.2.0)"},
		{"=1.2.4 || =1.3.0 && =1.2.0", "(=1.2.4 || (=1.3.0 && =1.2.0))"},
		{"(=1.2.4 || =1.3.0) && =1.2.0", "((=1.2.4 || =1.3.0) && =1.2.0)"},
		{"(=1.2.4 || !>1.3.0) && =1.2.0", "((=1.2.4 || !(>1.3.0)) && =1.2.0)"},
		{"!(=1.2.4 || >1.3.0) && =1.2.0", "(!(=1.2.4 || >1.3.0) && =1.2.0)"},
	}
	for i, test := range good {
		in := test.In
		out := test.Out
		t.Run(fmt.Sprintf("GoodString%03d", i), func(t *testing.T) {
			p, err := ParseConstraint(in)
			require.NoError(t, err, "error parsing %s", in)
			require.Equal(t, out, p.String())
			fmt.Printf("'%s' parsed as %s\n", in, p.String())
		})
	}

	bad := []string{
		"1.0.0",
		"= 1.0.0",
		">= 1.0.0",
		"> 1.0.0",
		"<= 1.0.0",
		"< 1.0.0",
		">>1.0.0",
		">1.0.0 =2.0.0",
		">1.0.0 &",
		"^1.1.1.1",
		"!1.0.0",
		">1.0.0 && 2.0.0",
		">1.0.0 | =2.0.0",
		"(>1.0.0 | =2.0.0)",
		"(>1.0.0 || =2.0.0",
		">1.0.0 || 2.0.0",
	}
	for i, s := range bad {
		in := s
		t.Run(fmt.Sprintf("BadString%03d", i), func(t *testing.T) {
			p, err := ParseConstraint(in)
			require.Nil(t, p)
			require.Error(t, err)
			fmt.Printf("'%s' parse error: %s\n", in, err)
		})
	}
}

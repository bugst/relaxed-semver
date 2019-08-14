//
// Copyright 2019 Cristian Maglie. All rights reserved.
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
}

func TestConstraintsParser(t *testing.T) {
	good := map[string]string{
		"":         "",
		"=1.3.0":   "=1.3.0",
		" =1.3.0 ": "=1.3.0",
		"=1.3.0 ":  "=1.3.0",
		" =1.3.0":  "=1.3.0",
		">=1.3.0":  ">=1.3.0",
		">1.3.0":   ">1.3.0",
		"<=1.3.0":  "<=1.3.0",
		"<1.3.0":   "<1.3.0",
	}
	for s, r := range good {
		p, err := ParseConstraint(s)
		require.NoError(t, err)
		require.Equal(t, r, p.String())
		fmt.Printf("'%s' parsed as %s\n", s, p.String())
	}
	bad := []string{
		"1.0.0",
		"= 1.0.0",
		">>1.0.0",
		">1.0.0 =2.0.0",
		">1.0.0 &",
	}
	for _, s := range bad {
		p, err := ParseConstraint(s)
		require.Nil(t, p)
		require.Error(t, err)
		fmt.Printf("'%s' parse error: %s\n", s, err)
	}
}

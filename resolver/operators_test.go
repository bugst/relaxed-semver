//
// Copyright 2019 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package resolver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConstraints(t *testing.T) {
	lt := &LessThanConstraint{v("1.3.0")}
	require.True(t, lt.Match(v("1.0.0")))
	require.False(t, lt.Match(v("1.3.0")))
	require.False(t, lt.Match(v("2.0.0")))
	require.Equal(t, "<1.3.0", lt.String())

	lte := &LessThanOrEqualConstraint{v("1.3.0")}
	require.True(t, lte.Match(v("1.0.0")))
	require.True(t, lte.Match(v("1.3.0")))
	require.False(t, lte.Match(v("2.0.0")))
	require.Equal(t, "<=1.3.0", lte.String())

	eq := &EqualsConstraint{v("1.3.0")}
	require.False(t, eq.Match(v("1.0.0")))
	require.True(t, eq.Match(v("1.3.0")))
	require.False(t, eq.Match(v("2.0.0")))
	require.Equal(t, "=1.3.0", eq.String())

	gte := &GreaterThanOrEqualConstraint{v("1.3.0")}
	require.False(t, gte.Match(v("1.0.0")))
	require.True(t, gte.Match(v("1.3.0")))
	require.True(t, gte.Match(v("2.0.0")))
	require.Equal(t, ">=1.3.0", gte.String())

	gt := &GreaterThanConstraint{v("1.3.0")}
	require.False(t, gt.Match(v("1.0.0")))
	require.False(t, gt.Match(v("1.3.0")))
	require.True(t, gt.Match(v("2.0.0")))
	require.Equal(t, ">1.3.0", gt.String())

	gt100 := &GreaterThanConstraint{v("1.0.0")}
	lte200 := &LessThanOrEqualConstraint{v("2.0.0")}
	and := &AndConstraint{[]Constraint{gt100, lte200}}
	require.False(t, and.Match(v("0.9.0")))
	require.False(t, and.Match(v("1.0.0")))
	require.True(t, and.Match(v("1.3.0")))
	require.True(t, and.Match(v("2.0.0")))
	require.False(t, and.Match(v("2.1.0")))
	require.Equal(t, "(>1.0.0 && <=2.0.0)", and.String())

	gt200 := &GreaterThanConstraint{v("2.0.0")}
	lte100 := &LessThanOrEqualConstraint{v("1.0.0")}
	or := &OrConstraint{[]Constraint{gt200, lte100}}
	require.True(t, or.Match(v("0.9.0")))
	require.True(t, or.Match(v("1.0.0")))
	require.False(t, or.Match(v("1.3.0")))
	require.False(t, or.Match(v("2.0.0")))
	require.True(t, or.Match(v("2.1.0")))
	require.Equal(t, "(>2.0.0 || <=1.0.0)", or.String())
}

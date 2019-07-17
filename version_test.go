//
// Copyright 2018 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func v(vers string) *Version {
	return MustParse(vers)
}

func TestVersionComparator(t *testing.T) {
	sign := map[int]string{1: ">", 0: "=", -1: "<"}
	ascending := func(list ...*Version) {
		for i := range list[0 : len(list)-1] {
			a := list[i]
			b := list[i+1]
			comp := a.CompareTo(b)
			fmt.Printf("%s %s %s\n", a, sign[comp], b)
			require.Equal(t, comp, -1)
			require.True(t, a.LessThan(b))
			require.True(t, a.LessThanOrEqual(b))
			require.False(t, a.Equal(b))
			require.False(t, a.GreaterThanOrEqual(b))
			require.False(t, a.GreaterThan(b))

			comp = b.CompareTo(a)
			fmt.Printf("%s %s %s\n", b, sign[comp], a)
			require.Equal(t, comp, 1)
			require.False(t, b.LessThan(a))
			require.False(t, b.LessThanOrEqual(a))
			require.False(t, b.Equal(a))
			require.True(t, b.GreaterThanOrEqual(a))
			require.True(t, b.GreaterThan(a))
		}
	}
	equal := func(list ...*Version) {
		for i, a := range list[:len(list)-1] {
			for _, b := range list[i+1:] {
				comp := a.CompareTo(b)
				fmt.Printf("%s %s %s\n", a, sign[comp], b)
				require.Equal(t, comp, 0)
				require.False(t, a.LessThan(b))
				require.True(t, a.LessThanOrEqual(b))
				require.True(t, a.Equal(b))
				require.True(t, a.GreaterThanOrEqual(b))
				require.False(t, a.GreaterThan(b))

				comp = b.CompareTo(a)
				fmt.Printf("%s %s %s\n", b, sign[comp], a)
				require.Equal(t, comp, 0)
				require.False(t, b.LessThan(a))
				require.True(t, b.LessThanOrEqual(a))
				require.True(t, b.Equal(a))
				require.True(t, b.GreaterThanOrEqual(a))
				require.False(t, b.GreaterThan(a))
			}
		}
	}
	ascending(
		MustParse("1.0.0-alpha"),
		MustParse("1.0.0-alpha.1"),
		MustParse("1.0.0-alpha.beta"),
		MustParse("1.0.0-beta"),
		MustParse("1.0.0-beta.2"),
		MustParse("1.0.0-beta.11"),
		MustParse("1.0.0-rc.1"),
		MustParse("1.0.0"),
		MustParse("1.0.1"),
		MustParse("1.1.1"),
		MustParse("1.6.22"),
		MustParse("1.8.1"),
		MustParse("2.1.1"),
	)
	equal(
		MustParse(""),
		MustParse("0"),
		MustParse("0.0"),
		MustParse("0.0.0"),
		MustParse("0+aaa"),
		MustParse("0.0+aaa"),
		MustParse("0.0.0+aaa"),
		MustParse("0+aaa.bbb"),
		MustParse("0.0+aaa.bbb"),
		MustParse("0.0.0+aaa.bbb"),
	)
	equal(
		MustParse("0-ab"),
		MustParse("0.0-ab"),
		MustParse("0.0.0-ab"),
		MustParse("0-ab+aaa"),
		MustParse("0.0-ab+aaa"),
		MustParse("0.0.0-ab+aaa"),
		MustParse("0-ab+aaa.bbb"),
		MustParse("0.0-ab+aaa.bbb"),
		MustParse("0.0.0-ab+aaa.bbb"),
	)
}

func TestNilVersionString(t *testing.T) {
	var nilVersion *Version
	require.Equal(t, "", nilVersion.String())
}

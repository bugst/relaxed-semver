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

var sign = map[int]string{1: ">", 0: "=", -1: "<"}

func v(vers string) *Version {
	return MustParse(vers)
}

func ascending(t *testing.T, allowEqual bool, list ...string) {
	for i := range list[0 : len(list)-1] {
		a := MustParse(list[i])
		b := MustParse(list[i+1])
		comp := a.CompareTo(b)
		fmt.Printf("%s %s %s\n", a, sign[comp], b)
		if allowEqual {
			require.LessOrEqual(t, comp, 0)
			require.True(t, a.LessThanOrEqual(b))
			require.False(t, a.GreaterThan(b))
		} else {
			require.Equal(t, comp, -1)
			require.True(t, a.LessThan(b))
			require.True(t, a.LessThanOrEqual(b))
			require.False(t, a.Equal(b))
			require.False(t, a.GreaterThanOrEqual(b))
			require.False(t, a.GreaterThan(b))
		}

		comp = b.CompareTo(a)
		fmt.Printf("%s %s %s\n", b, sign[comp], a)
		if allowEqual {
			require.GreaterOrEqual(t, comp, 0)
			require.False(t, b.LessThan(a))
			require.True(t, b.GreaterThanOrEqual(a))
		} else {
			require.Equal(t, comp, 1)
			require.False(t, b.LessThan(a))
			require.False(t, b.LessThanOrEqual(a))
			require.False(t, b.Equal(a))
			require.True(t, b.GreaterThanOrEqual(a))
			require.True(t, b.GreaterThan(a))
		}
	}
}

func TestVersionComparator(t *testing.T) {
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
	ascending(t, false,
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"1.0.0-alpha.beta",
		"1.0.0-beta",
		"1.0.0-beta.2",
		"1.0.0-beta.11",
		"1.0.0-rc.1",
		"1.0.0",
		"1.0.1",
		"1.1.1",
		"1.6.22",
		"1.8.1",
		"2.1.1",
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

func TestCompatibleWithVersionComparator(t *testing.T) {
	require.True(t, v("0.0.1-rc.0+build").CompatibleWith(v("0.0.1-rc.0")))
	list := []string{
		"0.0.1-rc.0",       // 0
		"0.0.1-rc.0+build", // 1
		"0.0.1-rc.1",       // 2
		"0.0.1",            // 3
		"0.0.1+build",      // 4
		"0.0.2-rc.1",       // 5 - BREAKING CHANGE
		"0.0.2-rc.1+build", // 6
		"0.0.2",            // 7
		"0.0.2+build",      // 8
		"0.0.3-rc.1",       // 9 - BREAKING CHANGE
		"0.0.3-rc.2",       // 10
		"0.0.3",            // 11
		"0.1.0",            // 12 - BREAKING CHANGE
		"0.3.3-rc.0",       // 13 - BREAKING CHANGE
		"0.3.3-rc.1",       // 14
		"0.3.3",            // 15
		"0.3.3+build",      // 16
		"0.3.4-rc.1",       // 17
		"0.3.4",            // 18
		"0.4.0",            // 19 - BREAKING CHANGE
		"1.0.0-rc",         // 20 - BREAKING CHANGE
		"1.0.0",            // 21
		"1.0.0+build",      // 22
		"1.2.1-rc",         // 23
		"1.2.1",            // 24
		"1.2.1+build",      // 25
		"1.2.3-rc.2",       // 26
		"1.2.3-rc.2+build", // 27
		"1.2.3",            // 28
		"1.2.3+build",      // 29
		"1.2.4",            // 30
		"1.3.0-rc.0+build", // 31
		"1.3.0",            // 32
		"1.3.0+build",      // 33
		"1.3.1-rc.0",       // 34
		"1.3.1-rc.1",       // 35
		"1.3.1",            // 36
		"1.3.5",            // 37
		"2.0.0-rc",         // 38 - BREAKING CHANGE
		"2.0.0-rc+build",   // 39
		"2.0.0",            // 40
		"2.0.0+build",      // 41
		"2.1.0-rc",         // 42
		"2.1.0-rc+build",   // 43
		"2.1.0",            // 44
		"2.1.0+build",      // 45
		"2.1.3-rc",         // 46
		"2.1.3",            // 47
		"2.3.0",            // 48
		"2.3.1",            // 49
		"3.0.0",            // 50 - BREAKING CHANGE
	}
	breaking := []int{5, 9, 12, 13, 19, 20, 38, 50}

	ascending(t, true, list...)
	compatible := func(which, from, to int) {
		x := MustParse(list[which])
		for _, comp := range list[:from] {
			y := MustParse(comp)
			require.False(t, x.CompatibleWith(y), "%s is not compatible with %s", x, y)
		}
		for _, comp := range list[from:to] {
			y := MustParse(comp)
			require.True(t, x.CompatibleWith(y), "%s is compatible with %s", x, y)
		}
		for _, comp := range list[to:] {
			y := MustParse(comp)
			require.False(t, x.CompatibleWith(y), "%s is not compatible with %s", x, y)
		}
	}

	j := 0
	for i := 0; i < len(list)-1; i++ {
		breakingIdx := 0
		for _, b := range breaking {
			if b > i {
				breakingIdx = b
				break
			}
		}
		if !MustParse(list[j]).Equal(MustParse(list[i])) {
			j = i
		}
		compatible(i, j, breakingIdx)
	}
}

func TestNilVersionString(t *testing.T) {
	var nilVersion *Version
	require.Equal(t, "", nilVersion.String())
}

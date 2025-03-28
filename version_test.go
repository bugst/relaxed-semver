//
// Copyright 2018-2025 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"cmp"
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
		if allowEqual {
			fmt.Printf("%s %s= %s\n", list[i], sign[comp], list[i+1])
			require.LessOrEqual(t, comp, 0)
			require.True(t, a.LessThanOrEqual(b))
			require.False(t, a.GreaterThan(b))
		} else {
			fmt.Printf("%s %s %s\n", list[i], sign[comp], list[i+1])
			require.Equal(t, comp, -1, "cmp(%s, %s) must return '<', but returned '%s'", list[i], list[i+1], sign[comp])
			require.True(t, a.LessThan(b))
			require.True(t, a.LessThanOrEqual(b))
			require.False(t, a.Equal(b))
			require.False(t, a.GreaterThanOrEqual(b))
			require.False(t, a.GreaterThan(b))
		}

		comp = b.CompareTo(a)
		fmt.Printf("%s %s %s\n", b, sign[comp], a)
		if allowEqual {
			require.GreaterOrEqual(t, comp, 0, "cmp(%s, %s) must return '>=', but returned '%s'", b, a, sign[comp])
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

	for i := range list[0 : len(list)-1] {
		a := MustParse(list[i]).SortableString()
		b := MustParse(list[i+1]).SortableString()
		comp := cmp.Compare(a, b)
		if allowEqual {
			fmt.Printf("%s %s= %s\n", list[i], sign[comp], list[i+1])
			require.LessOrEqual(t, comp, 0)
			require.True(t, a <= b)
			require.False(t, a > b)
		} else {
			fmt.Printf("%s %s %s\n", list[i], sign[comp], list[i+1])
			require.Equal(t, comp, -1, "cmp(%s, %s) (%s, %s) must return '<', but returned '%s'", list[i], list[i+1], a, b, sign[comp])
			require.True(t, a < b)
			require.True(t, a <= b)
			require.False(t, a == b)
			require.False(t, a >= b)
			require.False(t, a > b)
		}

		comp = cmp.Compare(b, a)
		fmt.Printf("%s %s %s\n", b, sign[comp], a)
		if allowEqual {
			require.GreaterOrEqual(t, comp, 0, "cmp(%s, %s) must return '>=', but returned '%s'", b, a, sign[comp])
			require.False(t, b < a)
			require.True(t, b >= a)
		} else {
			require.Equal(t, comp, 1)
			require.False(t, b < a)
			require.False(t, b <= a)
			require.False(t, b == a)
			require.True(t, b >= a)
			require.True(t, b > a)
		}
	}
}

func TestVersionComparator(t *testing.T) {
	equal := func(list ...*Version) {
		for i, a := range list[:len(list)-1] {
			for _, b := range list[i+1:] {
				comp := a.CompareTo(b)
				fmt.Printf("%s %s %s\n", a, sign[comp], b)
				require.Equal(t, comp, 0, "cmp(%s, %s) must return '=', but returned '%s'", a, b, sign[comp])
				require.False(t, a.LessThan(b), "NOT wanted: %s < %s", a, b)
				require.True(t, a.LessThanOrEqual(b), "wanted: %s <= %s", a, b)
				require.True(t, a.Equal(b), "wanted: %s = %s", a, b)
				require.True(t, a.GreaterThanOrEqual(b), "wanted: %s >= %s", a, b)
				require.False(t, a.GreaterThan(b), "NOT wanted: %s > %s", a, b)

				comp = b.CompareTo(a)
				fmt.Printf("%s %s %s\n", b, sign[comp], a)
				require.Equal(t, comp, 0, "cmp(%s, %s) must return '=', but returned '%s'", b, a, sign[comp])
				require.False(t, b.LessThan(a), "NOT wanted: %s < %s", b, a)
				require.True(t, b.LessThanOrEqual(a), "wanted: %s <= %s", b, a)
				require.True(t, b.Equal(a), "wanted: %s = %s", b, a)
				require.True(t, b.GreaterThanOrEqual(a), "wanted: %s >= %s", b, a)
				require.False(t, b.GreaterThan(a), "NOT wanted: %s > %s", b, a)
			}
		}
	}
	ascending(t, false, "", "0.0.1")
	ascending(t, false, "", "0.1")
	ascending(t, false, "", "1")
	ascending(t, false, "0", "0.0.1")
	ascending(t, false, "0", "0.1")
	ascending(t, false, "0", "1")
	ascending(t, false, "0.0", "0.0.1")
	ascending(t, false, "0.0", "0.1")
	ascending(t, false, "0.0", "1")
	ascending(t, false,
		"",
		"0.0.1",
		"0.1",
		"1.0.0-2",
		"1.0.0-11",
		"1.0.0-11a",
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"1.0.0-alpha.beta",
		"1.0.0-beta",
		"1.0.0-beta.2",
		"1.0.0-beta.11",
		"1.0.0-beta.11a",
		"1.0.0-rc.1",
		"1.0.0",
		"1.0.1",
		"1.1.1",
		"1.1.8",
		"1.1.22",
		"1.6.22",
		"1.8.1",
		"1.20.0",
		"2.1.1",
		"10.0.0",
		"17.3.0-atmel3.6.1-arduino7",
		"17.3.0-atmel3.6.1-arduino7not",
		"17.3.0-atmel3.6.1-beduino8",
		"17.3.0-atmel3.6.1-beduino8not",
		"17.3.0-atmel3a.6.1-arduino7",
		"17.3.0-atmel3a.16.2.arduino7",
		"17.3.0-atmel3a.16.12.arduino7",
		"17.3.0-atmel3a.16.1-arduino7",
		"17.3.0-atmel3a.16.12-arduino7",
		"17.3.0-atmel3a.16.2-arduino7",
		"34.0.0",
		"51.0.0",
		"99.0.0",
		"123.0.0",
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

func TestCompareNumbers(t *testing.T) {
	// ==
	require.Zero(t, compareNumber([]byte("0"), []byte("0")))
	require.Zero(t, compareNumber([]byte("5"), []byte("5")))
	require.Zero(t, compareNumber([]byte("15"), []byte("15")))

	// >
	testGreater := func(a, b string) {
		require.Positive(t, compareNumber([]byte(a), []byte(b)), `compareNumber("%s","%s") is not positive`, a, b)
		require.Negative(t, compareNumber([]byte(b), []byte(a)), `compareNumber("%s","%s") is not negative`, b, a)
	}
	testGreater("1", "")
	testGreater("1", "0")
	testGreater("1", "")
	testGreater("2", "1")
	testGreater("10", "")
	testGreater("10", "0")
	testGreater("10", "1")
	testGreater("10", "2")
}

func TestVersionGetters(t *testing.T) {
	type test struct {
		version    string
		prerelease string
		build      string
	}
	tests := []test{
		{"", "", ""},
		{"0", "", ""},
		{"1", "", ""},
		{"0.1", "", ""},
		{"1.1", "", ""},
		{"0.2.3", "", ""},
		{"1.2.3-aaa", "aaa", ""},
		{"0.2-aaa", "aaa", ""},
		{"1-aaa", "aaa", ""},
		{"0.2.3+bbb", "", "bbb"},
		{"1.3+bbb", "", "bbb"},
		{"0+bbb", "", "bbb"},
		{"1.2.3-aaa+bbb", "aaa", "bbb"},
		{"0.2-aaa+bbb", "aaa", "bbb"},
		{"1-aaa+bbb", "aaa", "bbb"},
		{"0.2.3-aaa.4.5.6+bbb.7.8.9", "aaa.4.5.6", "bbb.7.8.9"},
	}
	for _, tt := range tests {
		v := MustParse(tt.version)
		require.Equal(t, tt.version, v.String())
		require.Equal(t, tt.prerelease != "", v.IsPrerelease())
		require.Equal(t, tt.prerelease, v.Prerelease())
		require.Equal(t, tt.build != "", v.HasBuildMetadata())
		require.Equal(t, tt.build, v.BuildMetadata())
		r := ParseRelaxed(tt.version)
		require.Equal(t, tt.version, r.String())
		require.Equal(t, tt.prerelease != "", r.IsPrerelease())
		require.Equal(t, tt.prerelease, r.Prerelease())
		require.Equal(t, tt.build != "", r.HasBuildMetadata())
		require.Equal(t, tt.build, r.BuildMetadata())
	}
	relaxedTests := []test{
		{"asd", "", ""},
		{"123.123.123.123-123", "", ""},
		{"1.2.3-a@very@fancy@version", "", ""},
	}
	for _, tt := range relaxedTests {
		v, err := Parse(tt.version)
		require.Error(t, err, "should not parse %s", tt.version)
		require.Nil(t, v)
		r := ParseRelaxed(tt.version)
		require.Equal(t, tt.version, r.String())
		require.Equal(t, tt.prerelease != "", r.IsPrerelease())
		require.Equal(t, tt.prerelease, r.Prerelease())
		require.Equal(t, tt.build != "", r.HasBuildMetadata())
		require.Equal(t, tt.build, r.BuildMetadata())
	}
}

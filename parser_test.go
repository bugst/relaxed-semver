//
// Copyright 2018-2023 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	valid := func(in, normalized, expectedDump string) {
		v, err := Parse(in)
		require.NoError(t, err, "parsing '%s'", in)
		require.Equal(t, in, v.String(), "printing of '%s'", in)
		require.Equal(t, normalized, string(v.NormalizedString()), "normalized printing of '%s'", in)
		dump := string(v.major) + ","
		dump += string(v.minor) + ","
		dump += string(v.patch)
		for i, prerelease := range v.prerelases {
			dump += ","
			if v.numericPrereleases[i] {
				dump += "N"
			}
			dump += string(prerelease)
		}
		for _, build := range v.builds {
			dump += ";" + string(build)
		}
		require.Equal(t, expectedDump, dump, "fields of parsed '%s'", in)
		fmt.Printf("%s -> %s\n", in, v.String())
		v.Normalize()
		require.Equal(t, normalized, v.String(), "normalization of '%s'", in)
	}
	invalid := func(in string) {
		v, err := Parse(in)
		require.Error(t, err, "parsing '%s'", in)
		require.Nil(t, v, "parsed '%s'", in)
		fmt.Printf("%s -> %s\n", in, err)
	}

	valid("", "0.0.0", ",,")
	invalid("0.0.0.0")
	invalid("a")
	invalid(".")
	invalid("-ab")
	invalid("+ab")
	valid("0", "0.0.0", "0,,")
	valid("0.0.0", "0.0.0", "0,0,0")
	valid("1", "1.0.0", "1,,")
	valid("1.0.0", "1.0.0", "1,0,0")
	valid("14", "14.0.0", "14,,")
	valid("123456789123456789123456789", "123456789123456789123456789.0.0", "123456789123456789123456789,,")
	invalid("12ab")
	invalid("01")
	invalid("0ab")
	invalid(".1.1")
	invalid("1-")
	valid("1-0", "1.0.0-0", "1,,,N0")
	valid("1-pre", "1.0.0-pre", "1,,,pre")
	valid("1-pre.a", "1.0.0-pre.a", "1,,,pre,a")
	valid("1-pre.a.0", "1.0.0-pre.a.0", "1,,,pre,a,N0")
	valid("1-pre.0.a", "1.0.0-pre.0.a", "1,,,pre,N0,a")
	valid("1-pre.a.10", "1.0.0-pre.a.10", "1,,,pre,a,N10")
	invalid("1-pre.a.01")
	invalid("1-pre.a..1")
	invalid("1-pre.a.01.1")
	invalid("1-pre.a.01*.1")
	valid("1+build3", "1.0.0+build3", "1,,;build3")
	invalid("1+build3+build2")
	valid("1+build3.123.001", "1.0.0+build3.123.001", "1,,;build3;123;001")
	invalid("1+build3.123..001")
	invalid("1+build3.123*.001")
	valid("1-0+build3", "1.0.0-0+build3", "1,,,N0;build3")
	valid("1-pre+build3", "1.0.0-pre+build3", "1,,,pre;build3")
	valid("1-pre.a+build3", "1.0.0-pre.a+build3", "1,,,pre,a;build3")
	valid("1-pre.a.10+build3", "1.0.0-pre.a.10+build3", "1,,,pre,a,N10;build3")
	invalid("1-pre.a.01+build3")
	invalid("1-pre.a..1+build3")
	invalid("1-pre.a.01.1+build3")
	invalid("1-pre.a.01*.1+build3")
	valid("1-0+build3.123.001", "1.0.0-0+build3.123.001", "1,,,N0;build3;123;001")
	valid("1-pre+build3.123.001", "1.0.0-pre+build3.123.001", "1,,,pre;build3;123;001")
	valid("1-pre.a+build3.123.001", "1.0.0-pre.a+build3.123.001", "1,,,pre,a;build3;123;001")
	valid("1-pre.a.0+build3.123.001", "1.0.0-pre.a.0+build3.123.001", "1,,,pre,a,N0;build3;123;001")
	valid("1-pre.0.a+build3.123.001", "1.0.0-pre.0.a+build3.123.001", "1,,,pre,N0,a;build3;123;001")
	valid("1-pre.a.10+build3.123.001", "1.0.0-pre.a.10+build3.123.001", "1,,,pre,a,N10;build3;123;001")
	invalid("1-pre.a.+build3.123.001")
	invalid("1-pre.a.01+build3.123.001")
	invalid("1-pre.a.01*+build3.123.001")

	invalid("1.")
	invalid("1.a")
	invalid("1..2")
	valid("1.2", "1.2.0", "1,2,")
	valid("1.0", "1.0.0", "1,0,")
	invalid("1.02")
	invalid("1.0ab")
	invalid("1.12ab")
	valid("1.123456789123456789123456789", "1.123456789123456789123456789.0", "1,123456789123456789123456789,")
	invalid("1.2-")
	valid("1.2-0", "1.2.0-0", "1,2,,N0")
	valid("1.2-pre", "1.2.0-pre", "1,2,,pre")
	valid("1.2-pre.a", "1.2.0-pre.a", "1,2,,pre,a")
	valid("1.2-pre.a.0", "1.2.0-pre.a.0", "1,2,,pre,a,N0")
	valid("1.2-pre.0.a", "1.2.0-pre.0.a", "1,2,,pre,N0,a")
	valid("1.2-pre.a.10", "1.2.0-pre.a.10", "1,2,,pre,a,N10")
	valid("1.2-pre.a.10", "1.2.0-pre.a.10", "1,2,,pre,a,N10")
	invalid("1.2-pre.a.01")
	invalid("1.2-pre.a..1")
	invalid("1.2-pre.a.01.1")
	invalid("1.2-pre.a.01*.1")

	valid("1.2+build3", "1.2.0+build3", "1,2,;build3")
	valid("1.2-0+build3", "1.2.0-0+build3", "1,2,,N0;build3")
	invalid("1.2+build3+build2")
	valid("1.2+build3.123.001", "1.2.0+build3.123.001", "1,2,;build3;123;001")
	invalid("1.2+build3.123..001")
	invalid("1.2+build3.123*.001")
	valid("1.2-pre+build3", "1.2.0-pre+build3", "1,2,,pre;build3")
	valid("1.2-pre.a.0+build3", "1.2.0-pre.a.0+build3", "1,2,,pre,a,N0;build3")
	valid("1.2-pre.0.a+build3", "1.2.0-pre.0.a+build3", "1,2,,pre,N0,a;build3")
	valid("1.2-pre.a.10+build3", "1.2.0-pre.a.10+build3", "1,2,,pre,a,N10;build3")
	valid("1.2-pre.a+build3", "1.2.0-pre.a+build3", "1,2,,pre,a;build3")
	valid("1.2-pre.a.10+build3", "1.2.0-pre.a.10+build3", "1,2,,pre,a,N10;build3")
	invalid("1.2-pre.a.01+build3")
	invalid("1.2-pre.a..1+build3")
	invalid("1.2-pre.a.01.1+build3")
	invalid("1.2-pre.a.01*.1+build3")
	valid("1.2-0+build3.123.001", "1.2.0-0+build3.123.001", "1,2,,N0;build3;123;001")
	valid("1.2-pre+build3.123.001", "1.2.0-pre+build3.123.001", "1,2,,pre;build3;123;001")
	valid("1.2-pre.a+build3.123.001", "1.2.0-pre.a+build3.123.001", "1,2,,pre,a;build3;123;001")
	valid("1.2-pre.a.0+build3.123.001", "1.2.0-pre.a.0+build3.123.001", "1,2,,pre,a,N0;build3;123;001")
	valid("1.2-pre.0.a+build3.123.001", "1.2.0-pre.0.a+build3.123.001", "1,2,,pre,N0,a;build3;123;001")
	valid("1.2-pre.a.10+build3.123.001", "1.2.0-pre.a.10+build3.123.001", "1,2,,pre,a,N10;build3;123;001")
	valid("1.2-pre.a.10+build3.123.001", "1.2.0-pre.a.10+build3.123.001", "1,2,,pre,a,N10;build3;123;001")
	invalid("1.2-pre.a.+build3.123.001")
	invalid("1.2-pre.a.01+build3.123.001")
	invalid("1.2-pre.a.01*+build3.123.001")

	invalid("1.2.a")
	invalid("1.2.")
	valid("1.2.3", "1.2.3", "1,2,3")
	valid("1.2.0", "1.2.0", "1,2,0")
	invalid("1.2.03")
	invalid("1.2.0ab")
	invalid("1.2.34ab")
	valid("1.2.123456789123456789123456789", "1.2.123456789123456789123456789", "1,2,123456789123456789123456789")
	invalid("1.2.3-")
	valid("1.2.3-0", "1.2.3-0", "1,2,3,N0")
	valid("1.2.3-pre", "1.2.3-pre", "1,2,3,pre")
	valid("1.2.3-pre.a", "1.2.3-pre.a", "1,2,3,pre,a")
	valid("1.2.3-pre.a.0", "1.2.3-pre.a.0", "1,2,3,pre,a,N0")
	valid("1.2.3-pre.0.a", "1.2.3-pre.0.a", "1,2,3,pre,N0,a")
	valid("1.2.3-pre.a.10", "1.2.3-pre.a.10", "1,2,3,pre,a,N10")
	valid("1.2.3-pre.a.10", "1.2.3-pre.a.10", "1,2,3,pre,a,N10")
	invalid("1.2.3-pre.a.01")
	invalid("1.2.3-pre.a..1")
	invalid("1.2.3-pre.a.01.1")
	invalid("1.2.3-pre.a.01*.1")

	valid("1.2.3+build3", "1.2.3+build3", "1,2,3;build3")
	invalid("1.2.3+build3+build2")
	valid("1.2.3+build3.123.001", "1.2.3+build3.123.001", "1,2,3;build3;123;001")
	invalid("1.2.3+build3.123..001")
	invalid("1.2.3+build3.123*.001")
	valid("1.2.3-0+build3", "1.2.3-0+build3", "1,2,3,N0;build3")
	valid("1.2.3-pre+build3", "1.2.3-pre+build3", "1,2,3,pre;build3")
	valid("1.2.3-pre.a+build3", "1.2.3-pre.a+build3", "1,2,3,pre,a;build3")
	valid("1.2.3-pre.a.0+build3", "1.2.3-pre.a.0+build3", "1,2,3,pre,a,N0;build3")
	valid("1.2.3-pre.0.a+build3", "1.2.3-pre.0.a+build3", "1,2,3,pre,N0,a;build3")
	valid("1.2.3-pre.a.10+build3", "1.2.3-pre.a.10+build3", "1,2,3,pre,a,N10;build3")
	valid("1.2.3-pre.a.10+build3", "1.2.3-pre.a.10+build3", "1,2,3,pre,a,N10;build3")
	invalid("1.2.3-pre.a.01+build3")
	invalid("1.2.3-pre.a..1+build3")
	invalid("1.2.3-pre.a.01.1+build3")
	invalid("1.2.3-pre.a.01*.1+build3")
	valid("1.2.3-0+build3.123.001", "1.2.3-0+build3.123.001", "1,2,3,N0;build3;123;001")
	valid("1.2.3-pre+build3.123.001", "1.2.3-pre+build3.123.001", "1,2,3,pre;build3;123;001")
	valid("1.2.3-pre.a+build3.123.001", "1.2.3-pre.a+build3.123.001", "1,2,3,pre,a;build3;123;001")
	valid("1.2.3-pre.a.10+build3.123.001", "1.2.3-pre.a.10+build3.123.001", "1,2,3,pre,a,N10;build3;123;001")
	invalid("1.2.3-pre.a.+build3.123.001")
	invalid("1.2.3-pre.a.01+build3.123.001")
	invalid("1.2.3-pre.a.01*+build3.123.001")

	valid("1.2.3-pre.a-10.20.c-30", "1.2.3-pre.a-10.20.c-30", "1,2,3,pre,a-10,N20,c-30")
	valid("1.2.3--1-.23.1", "1.2.3--1-.23.1", "1,2,3,-1-,N23,N1")

	invalid("1.2.3.4")
	invalid("1.2.3.")
}

func TestNilVersionStringOutput(t *testing.T) {
	var nilVersion *Version
	require.Equal(t, "", nilVersion.String())
	require.Equal(t, "", string(nilVersion.NormalizedString()))
}

func TestParseRelaxed(t *testing.T) {
	bad := ParseRelaxed("bad")
	require.Nil(t, bad.version)
	require.Equal(t, []byte("bad"), bad.customversion)
	require.Equal(t, "bad", bad.String())
	good := ParseRelaxed("1.2.3-pre.a.10+build3.123.001")
	require.Nil(t, good.customversion)
	require.Equal(t, "1.2.3-pre.a.10+build3.123.001", good.version.String())
	require.Equal(t, "1.2.3-pre.a.10+build3.123.001", good.String())

}

func ExampleParseRelaxed() {
	WarnInvalidVersionWhenParsingRelaxed = true
	ParseRelaxed("bad")
	WarnInvalidVersionWhenParsingRelaxed = false

	// Output:
	// WARNING invalid semver version bad: no major version found
}

func TestMustParse(t *testing.T) {
	require.NotPanics(t, func() { MustParse("1.2.3") })
	require.Panics(t, func() { MustParse("bad") })
}

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

func TestRelaxedVersionComparator(t *testing.T) {
	sign := map[int]string{1: ">", 0: "=", -1: "<"}
	ascending := func(list ...*RelaxedVersion) {
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
	equal := func(list ...*RelaxedVersion) {
		for _, a := range list {
			for _, b := range list {
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
		ParseRelaxed("6_2"),
		ParseRelaxed("alpha"),
		ParseRelaxed("beta"),
		ParseRelaxed("gamma"),
		ParseRelaxed("1.0.0-alpha"),
		ParseRelaxed("1.0.0-alpha.1"),
		ParseRelaxed("1.0.0-alpha.beta"),
		ParseRelaxed("1.0.0-beta"),
		ParseRelaxed("1.0.0-beta.2"),
		ParseRelaxed("1.0.0-beta.11"),
		ParseRelaxed("1.0.0-rc.1"),
		ParseRelaxed("1.0.0"),
		ParseRelaxed("1.0.1"),
		ParseRelaxed("1.1.1"),
		ParseRelaxed("2.1.1"),
	)
	equal(
		ParseRelaxed(""),
		ParseRelaxed("0"),
		ParseRelaxed("0.0"),
		ParseRelaxed("0.0.0"),
		ParseRelaxed("0+aaa"),
		ParseRelaxed("0.0+aaa"),
		ParseRelaxed("0.0.0+aaa"),
		ParseRelaxed("0+aaa.bbb"),
		ParseRelaxed("0.0+aaa.bbb"),
		ParseRelaxed("0.0.0+aaa.bbb"),
	)
}

func TestRelaxedCompatibleWith(t *testing.T) {
	inv := ParseRelaxed("invalid-semver")
	inv2 := ParseRelaxed("invalid-semver-2")
	v145 := ParseRelaxed("1.4.5")
	v152 := ParseRelaxed("1.5.2")
	v213 := ParseRelaxed("2.1.3")
	require.True(t, inv.CompatibleWith(inv))
	require.False(t, inv.CompatibleWith(inv2))
	require.False(t, inv.CompatibleWith(v145))
	require.False(t, inv.CompatibleWith(v152))
	require.False(t, inv.CompatibleWith(v213))
	require.False(t, inv2.CompatibleWith(inv))
	require.True(t, inv2.CompatibleWith(inv2))
	require.False(t, inv2.CompatibleWith(v145))
	require.False(t, inv2.CompatibleWith(v152))
	require.False(t, inv2.CompatibleWith(v213))
	require.False(t, v145.CompatibleWith(inv))
	require.False(t, v145.CompatibleWith(inv2))
	require.True(t, v145.CompatibleWith(v145))
	require.True(t, v145.CompatibleWith(v152))
	require.False(t, v145.CompatibleWith(v213))
	require.False(t, v152.CompatibleWith(inv))
	require.False(t, v152.CompatibleWith(inv2))
	require.False(t, v152.CompatibleWith(v145))
	require.True(t, v152.CompatibleWith(v152))
	require.False(t, v152.CompatibleWith(v213))
	require.False(t, v213.CompatibleWith(inv))
	require.False(t, v213.CompatibleWith(inv2))
	require.False(t, v213.CompatibleWith(v145))
	require.False(t, v213.CompatibleWith(v152))
	require.True(t, v213.CompatibleWith(v213))
}

func TestNilRelaxedVersionString(t *testing.T) {
	var nilVersion *RelaxedVersion
	require.Equal(t, "", nilVersion.String())
	require.Equal(t, "", string(nilVersion.NormalizedString()))
	require.Equal(t, "1.0.0", string(ParseRelaxed("1.0.0").NormalizedString()))
	require.Equal(t, "invalid-semver", string(ParseRelaxed("invalid-semver").NormalizedString()))
}

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

func TestVersionComparator(t *testing.T) {
	sign := map[int]string{1: ">", 0: "=", -1: "<"}
	ascending := func(list ...*Version) {
		for i := range list[0 : len(list)-1] {
			a := list[i]
			b := list[i+1]
			comp := a.CompareTo(b)
			fmt.Printf("%s %s %s\n", a, sign[comp], b)
			require.Equal(t, comp, -1)
			comp = b.CompareTo(a)
			fmt.Printf("%s %s %s\n", b, sign[comp], a)
			require.Equal(t, comp, 1)
		}
	}
	equal := func(list ...*Version) {
		for _, a := range list {
			for _, b := range list {
				comp := a.CompareTo(b)
				fmt.Printf("%s %s %s\n", a, sign[comp], b)
				require.Equal(t, comp, 0)
				comp = b.CompareTo(a)
				fmt.Printf("%s %s %s\n", b, sign[comp], a)
				require.Equal(t, comp, 0)
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

//
// Copyright 2018 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListSorting(t *testing.T) {
	list := []*Version{
		MustParse("2.1.1"),
		MustParse("1.0.0-beta"),
		MustParse("1.0.0-beta.11"),
		MustParse("1.0.0-alpha.beta"),
		MustParse("1.0.0-alpha.1"),
		MustParse("1.0.0-rc.1"),
		MustParse("1.0.1"),
		MustParse("1.0.0"),
		MustParse("1.0.0-alpha"),
		MustParse("1.0.0-beta.2"),
		MustParse("1.1.1"),
	}
	ordered := []*Version{
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
	}
	sort.Sort(List(list))
	for i := range list {
		require.True(t, list[i].Equal(ordered[i]))
	}
}

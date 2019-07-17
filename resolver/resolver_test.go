//
// Copyright 2019 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package resolver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	semver "go.bug.st/relaxed-semver"
)

func v(vers string) *semver.Version {
	return semver.MustParse(vers)
}

func d(dep string) *Dependency {
	name := dep[0:1]
	switch dep[1:3] {
	case ">=":
		return &Dependency{Name: name, Constraint: &semver.GreaterThanOrEqual{Version: v(dep[3:])}}
	case "<=":
		return &Dependency{Name: name, Constraint: &semver.LessThanOrEqual{Version: v(dep[3:])}}
	}
	switch dep[1:2] {
	case "=":
		return &Dependency{Name: name, Constraint: &semver.Equals{Version: v(dep[2:])}}
	case ">":
		return &Dependency{Name: name, Constraint: &semver.GreaterThan{Version: v(dep[2:])}}
	case "<":
		return &Dependency{Name: name, Constraint: &semver.LessThan{Version: v(dep[2:])}}
	case "^":
		panic("'compatible with' operator not implemented: " + dep)
		// return &Dependency{Name: name, Constraint: &CompatibleWithConstraint{Version: v(dep[2:])}}
	default:
		panic("invalid operator in dep: " + dep)
	}
}

func deps(deps ...string) []*Dependency {
	res := []*Dependency{}
	for _, dep := range deps {
		res = append(res, d(dep))
	}
	return res
}

func rel(name, ver string, deps []*Dependency) *Release {
	return &Release{Name: name, Version: v(ver), Dependencies: deps}
}

func TestResolver(t *testing.T) {
	b131 := rel("B", "1.3.1", deps("C<2.0.0"))
	b130 := rel("B", "1.3.0", deps())
	b121 := rel("B", "1.2.1", deps())
	b120 := rel("B", "1.2.0", deps())
	b111 := rel("B", "1.1.1", deps())
	b110 := rel("B", "1.1.0", deps())
	b100 := rel("B", "1.0.0", deps())
	c200 := rel("C", "2.0.0", deps())
	c120 := rel("C", "1.2.0", deps())
	c111 := rel("C", "1.1.1", deps())
	c110 := rel("C", "1.1.0", deps())
	c102 := rel("C", "1.0.2", deps())
	c101 := rel("C", "1.0.1", deps())
	c100 := rel("C", "1.0.0", deps())
	c021 := rel("C", "0.2.1", deps())
	c020 := rel("C", "0.2.0", deps())
	c010 := rel("C", "0.1.0", deps())
	arch := &Archive{
		Releases: map[string]ReleasesSet{
			"B": ReleasesSet{b131, b130, b121, b120, b111, b110, b100},
			"C": ReleasesSet{c200, c120, c111, c110, c102, c101, c100, c021, c020, c010},
		},
	}

	a100 := rel("A", "1.0.0", deps("B>=1.2.0", "C>=2.0.0"))
	a110 := rel("A", "1.1.0", deps("B=1.2.0", "C>=2.0.0"))
	a120 := rel("A", "1.2.0", deps("B=1.2.0", "C>2.0.0"))

	verbose = true

	r1 := arch.Resolve(a100)
	require.Len(t, r1, 3)
	require.Contains(t, r1, a100)
	require.Contains(t, r1, b130)
	require.Contains(t, r1, c200)
	fmt.Println(r1)

	r2 := arch.Resolve(a110)
	require.Len(t, r2, 3)
	require.Contains(t, r2, a110)
	require.Contains(t, r2, b120)
	require.Contains(t, r2, c200)
	fmt.Println(r2)

	r3 := arch.Resolve(a120)
	require.Nil(t, r3)
	fmt.Println(r3)
}

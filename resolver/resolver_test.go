//
// Copyright 2019 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package resolver

import (
	"fmt"
	"testing"

	semver "go.bug.st/relaxed-semver"
)

func v(vers string) *semver.Version {
	return semver.MustParse(vers)
}

func d(dep string) *Dependency {
	name := dep[0:1]
	switch dep[1:3] {
	case ">=":
		return &Dependency{Name: name, Constraint: &GreaterThanOrEqualConstraint{Version: v(dep[3:])}}
	case "<=":
		return &Dependency{Name: name, Constraint: &LessThanOrEqualConstraint{Version: v(dep[3:])}}
	}
	switch dep[1:2] {
	case "=":
		return &Dependency{Name: name, Constraint: &EqualsConstraint{Version: v(dep[2:])}}
	case ">":
		return &Dependency{Name: name, Constraint: &GreaterThanConstraint{Version: v(dep[2:])}}
	case "<":
		return &Dependency{Name: name, Constraint: &LessThanConstraint{Version: v(dep[2:])}}
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
	arch := &Archive{
		Releases: map[string]ReleasesSet{
			"A": ReleasesSet{
				rel("A", "1.0.0", deps("B>=1.2.0", "C>=2.0.0")),
			},
			"B": ReleasesSet{
				rel("B", "1.3.1", deps("C<2.0.0")),
				rel("B", "1.3.0", deps()),
				rel("B", "1.2.1", deps()),
				rel("B", "1.2.0", deps()),
				rel("B", "1.1.1", deps()),
				rel("B", "1.1.0", deps()),
				rel("B", "1.0.0", deps()),
			},
			"C": ReleasesSet{
				rel("C", "2.0.0", deps()),
				rel("C", "1.2.0", deps()),
				rel("C", "1.1.1", deps()),
				rel("C", "1.1.0", deps()),
				rel("C", "1.0.2", deps()),
				rel("C", "1.0.1", deps()),
				rel("C", "1.0.0", deps()),
				rel("C", "0.2.1", deps()),
				rel("C", "0.2.0", deps()),
				rel("C", "0.1.0", deps()),
			},
		},
	}

	A := arch.Releases["A"]
	fmt.Println(arch.Resolve(A[0]))
}

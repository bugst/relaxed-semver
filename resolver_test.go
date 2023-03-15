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

type customDep struct {
	name string
	cond Constraint
}

// GetName return the name of the dependency (implements the Dependency interface)
func (c *customDep) GetName() string {
	return c.name
}

// GetConstraint return the version contraints of the dependency (implements the Dependency interface)
func (c *customDep) GetConstraint() Constraint {
	return c.cond
}

func (c *customDep) String() string {
	return c.name + c.cond.String()
}

type customRel struct {
	name string
	vers *Version
	deps []*customDep
}

// GetName return the name of the release (implements the Release interface)
func (r *customRel) GetName() string {
	return r.name
}

// GetVersion return the version of the release (implements the Release interface)
func (r *customRel) GetVersion() *Version {
	return r.vers
}

// GetDependencies return the dependencies of the release (implements the Release interface)
func (r *customRel) GetDependencies() []*customDep {
	return r.deps
}

func (r *customRel) String() string {
	return r.name + "@" + r.vers.String()
}

func d(t *testing.T, dep string) *customDep {
	name := dep[0:1]
	cond, err := ParseConstraint(dep[1:])
	require.NoError(t, err, "invalid operator in dep: %s (%s)", dep, err)
	return &customDep{name: name, cond: cond}
}

func deps(t *testing.T, deps ...string) []*customDep {
	res := []*customDep{}
	for _, dep := range deps {
		res = append(res, d(t, dep))
	}
	return res
}

func rel(name, ver string, deps []*customDep) *customRel {
	return &customRel{name: name, vers: v(ver), deps: deps}
}

func TestResolver(t *testing.T) {
	b131 := rel("B", "1.3.1", deps(t, "C<2.0.0"))
	b130 := rel("B", "1.3.0", deps(t))
	b121 := rel("B", "1.2.1", deps(t))
	b120 := rel("B", "1.2.0", deps(t))
	b111 := rel("B", "1.1.1", deps(t))
	b110 := rel("B", "1.1.0", deps(t))
	b100 := rel("B", "1.0.0", deps(t))
	c200 := rel("C", "2.0.0", deps(t))
	c120 := rel("C", "1.2.0", deps(t))
	c111 := rel("C", "1.1.1", deps(t, "B=1.1.1"))
	c110 := rel("C", "1.1.0", deps(t))
	c102 := rel("C", "1.0.2", deps(t))
	c101 := rel("C", "1.0.1", deps(t))
	c100 := rel("C", "1.0.0", deps(t))
	c021 := rel("C", "0.2.1", deps(t))
	c020 := rel("C", "0.2.0", deps(t))
	c010 := rel("C", "0.1.0", deps(t, "D"))
	d100 := rel("D", "1.0.0", deps(t))
	d120 := rel("D", "1.2.0", deps(t, "E"))
	e100 := rel("E", "1.0.0", deps(t))

	resolver := NewResolver[*customRel, *customDep]()
	resolver.AddReleases(Releases[*customRel, *customDep]{b131, b130, b121, b120, b111, b110, b100})
	resolver.AddReleases(Releases[*customRel, *customDep]{c200, c120, c111, c110, c102, c101, c100, c021, c020, c010})
	resolver.AddReleases(Releases[*customRel, *customDep]{d100, d120})
	resolver.AddReleases(Releases[*customRel, *customDep]{e100})

	a100 := rel("A", "1.0.0", deps(t, "B>=1.2.0", "C>=2.0.0"))
	a110 := rel("A", "1.1.0", deps(t, "B=1.2.0", "C>=2.0.0"))
	a111 := rel("A", "1.1.1", deps(t, "B", "C=1.1.1"))
	a120 := rel("A", "1.2.0", deps(t, "B=1.2.0", "C>2.0.0"))

	r1 := resolver.Resolve(a100, MaximumVersionAvailableStrategy)
	require.Len(t, r1, 3)
	require.Contains(t, r1, a100)
	require.Contains(t, r1, b130)
	require.Contains(t, r1, c200)
	fmt.Println(r1)

	r2 := resolver.Resolve(a110, MaximumVersionAvailableStrategy)
	require.Len(t, r2, 3)
	require.Contains(t, r2, a110)
	require.Contains(t, r2, b120)
	require.Contains(t, r2, c200)
	fmt.Println(r2)

	r3 := resolver.Resolve(a111, MaximumVersionAvailableStrategy)
	require.Len(t, r3, 3)
	require.Contains(t, r3, a111)
	require.Contains(t, r3, b111)
	require.Contains(t, r3, c111)
	fmt.Println(r3)

	r4 := resolver.Resolve(a120, MaximumVersionAvailableStrategy)
	require.Nil(t, r4)
	fmt.Println(r4)

	r5 := resolver.Resolve(c010, MaximumVersionAvailableStrategy)
	require.Contains(t, r5, c010)
	require.Contains(t, r5, d120)
	require.Contains(t, r5, e100)
	fmt.Println(r5)
}

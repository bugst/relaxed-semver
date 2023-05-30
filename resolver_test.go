//
// Copyright 2018-2022 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type customDep struct {
	name string
	cond Constraint
}

func (c *customDep) GetName() string {
	return c.name
}

func (c *customDep) GetConstraint() Constraint {
	return c.cond
}

func (c *customDep) String() string {
	return c.name + c.cond.String()
}

type customRel struct {
	name string
	vers *Version
	deps []Dependency
}

func (r *customRel) GetName() string {
	return r.name
}

func (r *customRel) GetVersion() *Version {
	return r.vers
}

func (r *customRel) GetDependencies() []Dependency {
	return r.deps
}

func (r *customRel) String() string {
	return r.name + "@" + r.vers.String()
}

func d(dep string) Dependency {
	name := dep[0:1]
	cond, err := ParseConstraint(dep[1:])
	if err != nil {
		panic("invalid operator in dep: " + dep + " (" + err.Error() + ")")
	}
	return &customDep{name: name, cond: cond}
}

func deps(deps ...string) []Dependency {
	res := []Dependency{}
	for _, dep := range deps {
		res = append(res, d(dep))
	}
	return res
}

func rel(name, ver string, deps []Dependency) Release {
	return &customRel{name: name, vers: v(ver), deps: deps}
}

func TestResolver(t *testing.T) {
	a100 := rel("A", "1.0.0", deps("B>=1.2.0", "C>=2.0.0"))
	a110 := rel("A", "1.1.0", deps("B=1.2.0", "C>=2.0.0"))
	a111 := rel("A", "1.1.1", deps("B", "C=1.1.1"))
	a120 := rel("A", "1.2.0", deps("B=1.2.0", "C>2.0.0"))
	a121 := rel("A", "1.2.1", deps("B", "C", "G", "H", "I", "E=1.0.1"))
	b131 := rel("B", "1.3.1", deps("C<2.0.0"))
	b130 := rel("B", "1.3.0", deps())
	b121 := rel("B", "1.2.1", deps())
	b120 := rel("B", "1.2.0", deps())
	b111 := rel("B", "1.1.1", deps())
	b110 := rel("B", "1.1.0", deps())
	b100 := rel("B", "1.0.0", deps())
	c200 := rel("C", "2.0.0", deps())
	c120 := rel("C", "1.2.0", deps())
	c111 := rel("C", "1.1.1", deps("B=1.1.1"))
	c110 := rel("C", "1.1.0", deps())
	c102 := rel("C", "1.0.2", deps())
	c101 := rel("C", "1.0.1", deps())
	c100 := rel("C", "1.0.0", deps())
	c021 := rel("C", "0.2.1", deps())
	c020 := rel("C", "0.2.0", deps())
	c010 := rel("C", "0.1.0", deps("D"))
	d100 := rel("D", "1.0.0", deps())
	d120 := rel("D", "1.2.0", deps("E"))
	e100 := rel("E", "1.0.0", deps())
	e101 := rel("E", "1.0.1", deps("F")) // INVALID
	g130 := rel("G", "1.3.0", deps())
	g140 := rel("G", "1.4.0", deps())
	g150 := rel("G", "1.5.0", deps())
	g160 := rel("G", "1.6.0", deps())
	g170 := rel("G", "1.7.0", deps())
	g180 := rel("G", "1.8.0", deps())
	h130 := rel("H", "1.3.0", deps())
	h140 := rel("H", "1.4.0", deps())
	h150 := rel("H", "1.5.0", deps())
	h160 := rel("H", "1.6.0", deps())
	h170 := rel("H", "1.7.0", deps())
	h180 := rel("H", "1.8.0", deps())
	i130 := rel("I", "1.3.0", deps())
	i140 := rel("I", "1.4.0", deps())
	i150 := rel("I", "1.5.0", deps())
	i160 := rel("I", "1.6.0", deps())
	i170 := rel("I", "1.7.0", deps())
	i180 := rel("I", "1.8.0", deps())
	arch := &Archive{
		Releases: map[string]Releases{
			"A": {a100, a110, a111, a120, a121},
			"B": {b131, b130, b121, b120, b111, b110, b100},
			"C": {c200, c120, c111, c110, c102, c101, c100, c021, c020, c010},
			"D": {d100, d120},
			"E": {e100, e101},
			"G": {g130, g140, g150, g160, g170, g180},
			"H": {h130, h140, h150, h160, h170, h180},
			"I": {i130, i140, i150, i160, i170, i180},
		},
	}

	a130 := rel("A", "1.3.0", deps())
	r0 := arch.Resolve(a130) // Non-existent in archive
	require.Nil(t, r0)

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

	r3 := arch.Resolve(a111)
	require.Len(t, r3, 3)
	require.Contains(t, r3, a111)
	require.Contains(t, r3, b111)
	require.Contains(t, r3, c111)
	fmt.Println(r3)

	r4 := arch.Resolve(a120)
	require.Nil(t, r4)
	fmt.Println(r4)

	r5 := arch.Resolve(c010)
	require.Contains(t, r5, c010)
	require.Contains(t, r5, d120)
	require.Contains(t, r5, e100)
	fmt.Println(r5)

	done := make(chan bool)
	go func() {
		r6 := arch.Resolve(a121)
		require.Nil(t, r6)
		fmt.Println(r6)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		require.FailNow(t, "test didn't complete in the allocated time")
	}
}

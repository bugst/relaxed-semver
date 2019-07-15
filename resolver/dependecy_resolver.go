//
// Copyright 2019 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package resolver

import (
	"fmt"

	semver "go.bug.st/relaxed-semver"
)

type Operator int

const (
	Equals             Operator = iota // =
	LessThan                           // <
	LessThanOrEqual                    // <=
	GreaterThan                        // >
	GreaterThanOrEqual                 // >=
	CompatibleWith                     // ^
)

func (o Operator) String() string {
	switch o {
	case Equals:
		return "="
	case LessThan:
		return "<"
	case LessThanOrEqual:
		return "<="
	case GreaterThan:
		return ">"
	case GreaterThanOrEqual:
		return ">="
	case CompatibleWith:
		return "^"
	}
	return "?"
}

type Dependency struct {
	Name    string
	Version *semver.Version
	Op      Operator
}

func (d *Dependency) String() string {
	return d.Name + d.Op.String() + d.Version.String()
}

type Release struct {
	Name         string
	Version      *semver.Version
	Dependencies []*Dependency
}

func (r *Release) String() string {
	return r.Name + "@" + r.Version.String()
}

func (r *Release) Match(dep *Dependency) bool {
	if r.Name != dep.Name {
		return false
	}
	switch dep.Op {
	case Equals:
		return r.Version.Equal(dep.Version)
	case LessThan:
		return r.Version.LessThan(dep.Version)
	case LessThanOrEqual:
		return r.Version.LessThanOrEqual(dep.Version)
	case GreaterThan:
		return r.Version.GreaterThan(dep.Version)
	case GreaterThanOrEqual:
		return r.Version.GreaterThanOrEqual(dep.Version)
	case CompatibleWith:
		panic("todo: " + dep.String())
	}
	panic("invalid dependency " + dep.String())
}

type ReleasesSet []*Release

func (set ReleasesSet) FilterBy(dep *Dependency) ReleasesSet {
	res := []*Release{}
	for _, release := range set {
		if release.Match(dep) {
			res = append(res, release)
		}
	}
	return res
}

type Archive struct {
	Releases map[string]ReleasesSet
}

func (ar *Archive) Resolve(release *Release) []*Release {
	solution := map[string]*Release{release.Name: release}
	depsToProcess := release.Dependencies
	return ar.resolve(solution, depsToProcess)
}

func (ar *Archive) resolve(solution map[string]*Release, depsToProcess []*Dependency) []*Release {
	debug := func(msg string) {
		for i := 0; i < len(solution); i++ {
			fmt.Print("   ")
		}
		fmt.Println(msg)
	}
	if len(depsToProcess) == 0 {
		debug("All dependencies have been resolved.")
		debug(fmt.Sprintf(">> %s", solution))
		res := []*Release{}
		for _, v := range solution {
			res = append(res, v)
		}
		return res
	}

	// Pick the first dependency in the deps to process
	dep := depsToProcess[0]
	debug(fmt.Sprintf("Considering next dep: %s", dep))

	// If a release is already picked in the solution check if it match the dep
	if existingRelease, has := solution[dep.Name]; has {
		debug("already in solution...")
		if existingRelease.Match(dep) {
			debug("...and the release match the dependency, go on")
			return ar.resolve(solution, depsToProcess[1:])
		}
		debug("...and the release do NOT match dependency, rollback")
		return nil
	}

	// Otherwise start backtracking the dependency
	releases := ar.Releases[dep.Name].FilterBy(dep)
	debug(fmt.Sprintf("releases matching criteria: %s", releases))
	for _, release := range releases {
		debug(fmt.Sprintf("try with %s", release))
		solution[dep.Name] = release
		res := ar.resolve(solution, append(depsToProcess[1:], release.Dependencies...))
		if res != nil {
			return res
		}
		debug(fmt.Sprintf("%s did not work...", release))
		delete(solution, dep.Name)
	}
	return nil
}

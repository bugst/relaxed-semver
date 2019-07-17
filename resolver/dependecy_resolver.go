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

// Dependency represents a dependency, it must provide methods to return Name and Constraints
type Dependency interface {
	Name() string
	Constraint() semver.Constraint
}

type Release struct {
	Name         string
	Version      *semver.Version
	Dependencies []Dependency
}

func (r *Release) String() string {
	return r.Name + "@" + r.Version.String()
}

func (r *Release) Match(dep Dependency) bool {
	return r.Name == dep.Name() && dep.Constraint().Match(r.Version)
}

type ReleasesSet []*Release

func (set ReleasesSet) FilterBy(dep Dependency) ReleasesSet {
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

// To be redefined in Tests to increase output
var verbose = false

func (ar *Archive) resolve(solution map[string]*Release, depsToProcess []Dependency) []*Release {
	debug := func(msg string) {}
	if verbose {
		debug = func(msg string) {
			for i := 0; i < len(solution); i++ {
				fmt.Print("   ")
			}
			fmt.Println(msg)
		}
	}
	debug(fmt.Sprintf("deps to process: %s", depsToProcess))
	if len(depsToProcess) == 0 {
		debug("All dependencies have been resolved.")
		res := []*Release{}
		for _, v := range solution {
			res = append(res, v)
		}
		return res
	}

	// Pick the first dependency in the deps to process
	dep := depsToProcess[0]
	depName := dep.Name()
	debug(fmt.Sprintf("Considering next dep: %s", dep))

	// If a release is already picked in the solution check if it match the dep
	if existingRelease, has := solution[depName]; has {
		if existingRelease.Match(dep) {
			debug(fmt.Sprintf("%s already in solution and matching", existingRelease))
			return ar.resolve(solution, depsToProcess[1:])
		}
		debug(fmt.Sprintf("%s already in solution do not match... rollingback", existingRelease))
		return nil
	}

	// Otherwise start backtracking the dependency
	releases := ar.Releases[dep.Name()].FilterBy(dep)
	debug(fmt.Sprintf("releases matching criteria: %s", releases))
	for _, release := range releases {
		debug(fmt.Sprintf("try with %s %s", release, release.Dependencies))
		solution[depName] = release
		res := ar.resolve(solution, append(depsToProcess[1:], release.Dependencies...))
		if res != nil {
			return res
		}
		debug(fmt.Sprintf("%s did not work...", release))
		delete(solution, depName)
	}
	return nil
}

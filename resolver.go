//
// Copyright 2018-2023 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import "sort"

// Dependency represents a dependency, it must provide methods to return Name and Constraints
type Dependency interface {
	GetName() string
	GetConstraint() Constraint
}

// Release represents a release, it must provide methods to return Name, Version and Dependencies
type Release interface {
	GetName() string
	GetVersion() *Version
	GetDependencies() []Dependency
}

func match(r Release, dep Dependency) bool {
	return r.GetName() == dep.GetName() && dep.GetConstraint().Match(r.GetVersion())
}

// Releases is a list of Release
type Releases []Release

// FilterBy return a subset of the Releases matching the provided Dependency
func (set Releases) FilterBy(dep Dependency) Releases {
	res := []Release{}
	for _, r := range set {
		if match(r, dep) {
			res = append(res, r)
		}
	}
	return res
}

// SortDescent sort the Releases in this set in descending order (the lastest
// release is the first)
func (set Releases) SortDescent() {
	sort.Slice(set, func(i, j int) bool {
		return set[i].GetVersion().GreaterThan(set[j].GetVersion())
	})
}

// Archive contains all Releases set to consider for dependency resolution
type Archive struct {
	Releases map[string]Releases
}

// Resolve will try to depp-resolve dependencies from the Release passed as
// arguent using a backtracking algorithm.
func (ar *Archive) Resolve(release Release) []Release {
	mainDep := &bareDependency{
		name:    release.GetName(),
		version: release.GetVersion(),
	}
	return ar.resolve(map[string]Release{}, []Dependency{mainDep}, map[Dependency]int{})
}

type bareDependency struct {
	name    string
	version *Version
}

func (b *bareDependency) GetName() string {
	return b.name
}

func (b *bareDependency) GetConstraint() Constraint {
	return &Equals{Version: b.version}
}

func (b *bareDependency) String() string {
	return b.GetName() + b.GetConstraint().String()
}

func (ar *Archive) resolve(solution map[string]Release, depsToProcess []Dependency, problematicDeps map[Dependency]int) []Release {
	debug("deps to process: %s", depsToProcess)
	if len(depsToProcess) == 0 {
		debug("All dependencies have been resolved.")
		res := []Release{}
		for _, v := range solution {
			res = append(res, v)
		}
		return res
	}

	// Pick the first dependency in the deps to process
	dep := depsToProcess[0]
	depName := dep.GetName()
	debug("Considering next dep: %s", depName)

	// If a release is already picked in the solution check if it match the dep
	if existingRelease, has := solution[depName]; has {
		if match(existingRelease, dep) {
			debug("%s already in solution and matching", existingRelease)
			return ar.resolve(solution, depsToProcess[1:], problematicDeps)
		}
		debug("%s already in solution do not match... rollingback", existingRelease)
		return nil
	}

	// Otherwise start backtracking the dependency
	releases := ar.Releases[dep.GetName()].FilterBy(dep)

	// Consider the latest versions first
	releases.SortDescent()

	findMissingDeps := func(deps []Dependency) Dependency {
		for _, dep := range deps {
			if _, ok := ar.Releases[dep.GetName()]; !ok {
				return dep
			}
		}
		return nil
	}

	debug("releases matching criteria: %s", releases)
	for _, release := range releases {
		deps := release.GetDependencies()
		debug("try with %s %s", release, deps)

		if missingDep := findMissingDeps(deps); missingDep != nil {
			debug("%s did not work, becuase his dependency %s does not exists", release, missingDep.GetName())
			continue
		}

		solution[depName] = release
		newDepsToProcess := append(depsToProcess[1:], deps...)
		// bubble up problematics deps so they are processed first
		sort.Slice(newDepsToProcess, func(i, j int) bool {
			return problematicDeps[newDepsToProcess[i]] > problematicDeps[newDepsToProcess[j]]
		})
		if res := ar.resolve(solution, newDepsToProcess, problematicDeps); res != nil {
			return res
		}
		debug("%s did not work...", release)
		delete(solution, depName)
	}

	problematicDeps[dep]++
	return nil
}

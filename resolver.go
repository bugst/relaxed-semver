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
type Release[D Dependency] interface {
	GetName() string
	GetVersion() *Version
	GetDependencies() []D
}

// Releases is a list of Release of the same package (all releases with
// the same Name but different Version)
type Releases[R Release[D], D Dependency] []R

// FilterBy return a subset of the Releases matching the provided Constraint
func (set Releases[R, D]) FilterBy(c Constraint) Releases[R, D] {
	var res Releases[R, D]
	for _, r := range set {
		if c.Match(r.GetVersion()) {
			res = append(res, r)
		}
	}
	return res
}

// SortDescent sort the Releases in this set in descending order (the lastest
// release is the first)
func (set Releases[R, D]) SortDescent() {
	sort.Slice(set, func(i, j int) bool {
		return set[i].GetVersion().GreaterThan(set[j].GetVersion())
	})
}

// Resolver is a container with references to all Releases to consider for
// dependency resolution
type Resolver[R Release[D], D Dependency] struct {
	releases map[string]Releases[R, D]
}

// NewResolver creates a new archive
func NewResolver[R Release[D], D Dependency]() *Resolver[R, D] {
	return &Resolver[R, D]{
		releases: map[string]Releases[R, D]{},
	}
}

// AddRelease adds a release to this archive
func (ar *Resolver[R, D]) AddRelease(rel R) {
	relName := rel.GetName()
	ar.releases[relName] = append(ar.releases[relName], rel)
}

// AddReleases adds all the releases to this archive
func (ar *Resolver[R, D]) AddReleases(rels ...R) {
	for _, rel := range rels {
		relName := rel.GetName()
		ar.releases[relName] = append(ar.releases[relName], rel)
	}
}

// Resolve will try to depp-resolve dependencies from the Release passed as
// arguent using a backtracking algorithm. This function is NOT thread-safe.
func (ar *Resolver[R, D]) Resolve(release R) Releases[R, D] {
	// Initial empty state of the resolver
	solution := map[string]R{}
	depsToProcess := []D{}
	problematicDeps := map[dependencyHash]int{}

	// Check if the release is in the archive
	if len(ar.releases[release.GetName()].FilterBy(&Equals{Version: release.GetVersion()})) == 0 {
		return nil
	}

	// Add the requested release to the solution and proceed
	// with the dependencies resolution
	solution[release.GetName()] = release
	depsToProcess = append(depsToProcess, release.GetDependencies()...)
	return ar.resolve(solution, depsToProcess, problematicDeps)
}

type dependencyHash string

func hashDependency[D Dependency](dep D) dependencyHash {
	return dependencyHash(dep.GetName() + "/" + dep.GetConstraint().String())
}

func (ar *Resolver[R, D]) resolve(solution map[string]R, depsToProcess []D, problematicDeps map[dependencyHash]int) Releases[R, D] {
	debug("deps to process: %s", depsToProcess)
	if len(depsToProcess) == 0 {
		debug("All dependencies have been resolved.")
		var res Releases[R, D]
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
		if dep.GetConstraint().Match(existingRelease.GetVersion()) {
			debug("%s already in solution and matching", existingRelease)
			return ar.resolve(solution, depsToProcess[1:], problematicDeps)
		}
		debug("%s already in solution do not match... rollingback", existingRelease)
		return nil
	}

	// Otherwise start backtracking the dependency
	releases := ar.releases[dep.GetName()].FilterBy(dep.GetConstraint())

	// Consider the latest versions first
	releases.SortDescent()

	debug("releases matching criteria: %s", releases)
backtracking_loop:
	for _, release := range releases {
		releaseDeps := release.GetDependencies()
		debug("try with %s %s", release, releaseDeps)

		for _, releaseDep := range releaseDeps {
			if _, ok := ar.releases[releaseDep.GetName()]; !ok {
				debug("%s did not work, becuase his dependency %s does not exists", release, releaseDep.GetName())
				continue backtracking_loop
			}
		}

		solution[depName] = release
		newDepsToProcess := append(depsToProcess[1:], deps...)
		// bubble up problematics deps so they are processed first
		sort.Slice(newDepsToProcess, func(i, j int) bool {
			ci := hashDependency(newDepsToProcess[i])
			cj := hashDependency(newDepsToProcess[j])
			return problematicDeps[ci] > problematicDeps[cj]
		})
		if res := ar.resolve(solution, newDepsToProcess, problematicDeps); res != nil {
			return res
		}
		debug("%s did not work...", release)
		delete(solution, depName)
	}

	problematicDeps[hashDependency(dep)]++
	return nil
}

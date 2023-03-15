//
// Copyright 2018-2022 Cristian Maglie. All rights reserved.
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

// Releases is an array of Release
type Releases[R Release[D], D Dependency] []R

func match[R Release[D], D Dependency](r R, dep D) bool {
	return r.GetName() == dep.GetName() && dep.GetConstraint().Match(r.GetVersion())
}

// FilterBy return a subset of the Releases matching the provided Dependency
func (set Releases[R, D]) FilterBy(dep D) Releases[R, D] {
	res := []R{}
	for _, r := range set {
		if match(r, dep) {
			res = append(res, r)
		}
	}
	return res
}

// SortDescending sort the Releases in this set in descending order (the lastest
// release is the first)
func (set Releases[R, D]) SortDescending() {
	sort.Slice(set, func(i, j int) bool {
		return set[i].GetVersion().GreaterThan(set[j].GetVersion())
	})
}

// SortAscending sort the Releases in this set in ascending order (the lastest
// release is the latest)
func (set Releases[R, D]) SortAscending() {
	sort.Slice(set, func(i, j int) bool {
		return set[i].GetVersion().LessThan(set[j].GetVersion())
	})
}

// ResolutionStrategy is a resolution strategy for the dependency resolver
type ResolutionStrategy int

// MinimumVersionRequiredStrategy is a resolution strategy where the minimum
// required version is taken from the available set. The computed solution
// is stable, running the algorithm again will produce the same solution even
// if new releases becomes available in the meantime.
const MinimumVersionRequiredStrategy ResolutionStrategy = iota

// MaximumVersionAvailableStrategy is a strategy where the latest available
// version is taken from the available set. The computed solution will have
// all the latest possible versions of the libraries. The solution is not
// stable, it will probably require a lock file to be able to reproduce it
// at a later time.
const MaximumVersionAvailableStrategy ResolutionStrategy = iota

// Resolve will try to resolve dependencies of the target Release using the
// given resolution strategy.
func (r *Releases[R, D]) Resolve(target R, strategy ResolutionStrategy) []R {
	index := map[string]Releases[R, D]{}
	for _, release := range *r {
		index[release.GetName()] = append(index[release.GetName()], release)
	}
	ar := Resolver[R, D]{
		index:         index,
		solution:      map[string]R{target.GetName(): target},
		depsToProcess: target.GetDependencies(),
		strategy:      strategy,
	}
	return ar.resolve()
}

// Resolver is a dependency resolver, it must be created with NewResolver method
type Resolver[R Release[D], D Dependency] struct {
	index         map[string]Releases[R, D]
	solution      map[string]R
	depsToProcess []D
	strategy      ResolutionStrategy
}

// NewResolver creates a new Resolver
func NewResolver[R Release[D], D Dependency]() *Resolver[R, D] {
	return &Resolver[R, D]{
		index: map[string]Releases[R, D]{},
	}
}

func (ar *Resolver[R, D]) AddReleases(releases Releases[R, D]) {
	for _, release := range releases {
		ar.AddRelease(release)
	}
}

func (ar *Resolver[R, D]) AddRelease(release R) {
	ar.index[release.GetName()] = append(ar.index[release.GetName()], release)
}

func (ar *Resolver[R, D]) Resolve(target R, strategy ResolutionStrategy) []R {
	ar.solution = map[string]R{target.GetName(): target}
	ar.depsToProcess = target.GetDependencies()
	ar.strategy = strategy
	return ar.resolve()
}

func (ar *Resolver[R, D]) resolve() []R {
	debug("deps to process: %s", ar.depsToProcess)
	if len(ar.depsToProcess) == 0 {
		debug("All dependencies have been resolved.")
		res := []R{}
		for _, v := range ar.solution {
			res = append(res, v)
		}
		return res
	}

	// Pick the first dependency in the deps to process
	dep := ar.depsToProcess[0]
	depName := dep.GetName()
	debug("Considering next dep: %s", dep)

	// If a release is already picked in the solution check if it match the dep
	if existingRelease, has := ar.solution[depName]; has {
		if match(existingRelease, dep) {
			debug("%s already in solution and matching", existingRelease)
			ar.depsToProcess = ar.depsToProcess[1:]
			return ar.resolve()
		}
		debug("%s already in solution do not match... rollingback", existingRelease)
		return nil
	}

	// Otherwise start backtracking the dependency
	releases := ar.index[dep.GetName()].FilterBy(dep)

	// Consider the latest versions first
	if ar.strategy == MaximumVersionAvailableStrategy {
		releases.SortDescending()
	} else {
		releases.SortAscending()
	}

	debug("releases matching criteria: %s", releases)
	for _, release := range releases {
		debug("try with %s %s", release, release.GetDependencies())
		ar.solution[depName] = release
		rolledBackDeps := ar.depsToProcess
		ar.depsToProcess = append(ar.depsToProcess[1:], release.GetDependencies()...)
		res := ar.resolve()
		if res != nil {
			return res
		}
		ar.depsToProcess = rolledBackDeps
		debug("%s did not work...", release)
		delete(ar.solution, depName)
	}
	return nil
}

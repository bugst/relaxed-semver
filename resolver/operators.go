//
// Copyright 2019 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package resolver

import (
	semver "go.bug.st/relaxed-semver"
)

// Constraint is a condition that a Version can match or not
type Constraint interface {
	// Match returns true if the Version satisfies the condition
	Match(*semver.Version) bool

	String() string
}

// EqualsConstraint is the equality (=) constraint
type EqualsConstraint struct {
	Version *semver.Version
}

// Match returns true if v satisfies the condition
func (eq *EqualsConstraint) Match(v *semver.Version) bool {
	return v.Equal(eq.Version)
}

func (eq *EqualsConstraint) String() string {
	return "=" + eq.Version.String()
}

// LessThanConstraint is the less than (<) constraint
type LessThanConstraint struct {
	Version *semver.Version
}

// Match returns true if v satisfies the condition
func (lt *LessThanConstraint) Match(v *semver.Version) bool {
	return v.LessThan(lt.Version)
}

func (lt *LessThanConstraint) String() string {
	return "<" + lt.Version.String()
}

// LessThanOrEqualConstraint is the "less than or equal" (<=) constraint
type LessThanOrEqualConstraint struct {
	Version *semver.Version
}

// Match returns true if v satisfies the condition
func (lte *LessThanOrEqualConstraint) Match(v *semver.Version) bool {
	return v.LessThanOrEqual(lte.Version)
}

func (lte *LessThanOrEqualConstraint) String() string {
	return "<=" + lte.Version.String()
}

// GreaterThanConstraint is the "greater than" (>) constraint
type GreaterThanConstraint struct {
	Version *semver.Version
}

// Match returns true if v satisfies the condition
func (gt *GreaterThanConstraint) Match(v *semver.Version) bool {
	return v.GreaterThan(gt.Version)
}

func (gt *GreaterThanConstraint) String() string {
	return ">" + gt.Version.String()
}

// GreaterThanOrEqualConstraint is the "greater than or equal" (>=) constraint
type GreaterThanOrEqualConstraint struct {
	Version *semver.Version
}

// Match returns true if v satisfies the condition
func (gte *GreaterThanOrEqualConstraint) Match(v *semver.Version) bool {
	return v.GreaterThanOrEqual(gte.Version)
}

func (gte *GreaterThanOrEqualConstraint) String() string {
	return ">=" + gte.Version.String()
}

// OrConstraint will match if ANY of the Operands Constraint will match
type OrConstraint struct {
	Operands []Constraint
}

// Match returns true if v satisfies the condition
func (or *OrConstraint) Match(v *semver.Version) bool {
	for _, op := range or.Operands {
		if op.Match(v) {
			return true
		}
	}
	return false
}

func (or *OrConstraint) String() string {
	res := "("
	for i, op := range or.Operands {
		if i > 0 {
			res += " || "
		}
		res += op.String()
	}
	res += ")"
	return res
}

// AndConstraint will match if ALL the Operands Constraint will match
type AndConstraint struct {
	Operands []Constraint
}

// Match returns true if v satisfies the condition
func (and *AndConstraint) Match(v *semver.Version) bool {
	for _, op := range and.Operands {
		if !op.Match(v) {
			return false
		}
	}
	return true
}

func (and *AndConstraint) String() string {
	res := "("
	for i, op := range and.Operands {
		if i > 0 {
			res += " && "
		}
		res += op.String()
	}
	res += ")"
	return res
}

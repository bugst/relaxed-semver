//
// Copyright 2019 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

// Constraint is a condition that a Version can match or not
type Constraint interface {
	// Match returns true if the Version satisfies the condition
	Match(*Version) bool

	String() string
}

// Equals is the equality (=) constraint
type Equals struct {
	Version *Version
}

// Match returns true if v satisfies the condition
func (eq *Equals) Match(v *Version) bool {
	return v.Equal(eq.Version)
}

func (eq *Equals) String() string {
	return "=" + eq.Version.String()
}

// LessThan is the less than (<) constraint
type LessThan struct {
	Version *Version
}

// Match returns true if v satisfies the condition
func (lt *LessThan) Match(v *Version) bool {
	return v.LessThan(lt.Version)
}

func (lt *LessThan) String() string {
	return "<" + lt.Version.String()
}

// LessThanOrEqual is the "less than or equal" (<=) constraint
type LessThanOrEqual struct {
	Version *Version
}

// Match returns true if v satisfies the condition
func (lte *LessThanOrEqual) Match(v *Version) bool {
	return v.LessThanOrEqual(lte.Version)
}

func (lte *LessThanOrEqual) String() string {
	return "<=" + lte.Version.String()
}

// GreaterThan is the "greater than" (>) constraint
type GreaterThan struct {
	Version *Version
}

// Match returns true if v satisfies the condition
func (gt *GreaterThan) Match(v *Version) bool {
	return v.GreaterThan(gt.Version)
}

func (gt *GreaterThan) String() string {
	return ">" + gt.Version.String()
}

// GreaterThanOrEqual is the "greater than or equal" (>=) constraint
type GreaterThanOrEqual struct {
	Version *Version
}

// Match returns true if v satisfies the condition
func (gte *GreaterThanOrEqual) Match(v *Version) bool {
	return v.GreaterThanOrEqual(gte.Version)
}

func (gte *GreaterThanOrEqual) String() string {
	return ">=" + gte.Version.String()
}

// Or will match if ANY of the Operands Constraint will match
type Or struct {
	Operands []Constraint
}

// Match returns true if v satisfies the condition
func (or *Or) Match(v *Version) bool {
	for _, op := range or.Operands {
		if op.Match(v) {
			return true
		}
	}
	return false
}

func (or *Or) String() string {
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

// And will match if ALL the Operands Constraint will match
type And struct {
	Operands []Constraint
}

// Match returns true if v satisfies the condition
func (and *And) Match(v *Version) bool {
	for _, op := range and.Operands {
		if !op.Match(v) {
			return false
		}
	}
	return true
}

func (and *And) String() string {
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

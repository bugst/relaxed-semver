//
// Copyright 2019 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"fmt"
	"strings"
)

// Constraint is a condition that a Version can match or not
type Constraint interface {
	// Match returns true if the Version satisfies the condition
	Match(*Version) bool

	String() string
}

// ParseConstraint converts a string into a Constraint. The resulting Constraint
// may be converted back to string using the String() method.
// WIP: only simple constraint (like ""=1.2.0" or ">=2.0.0) are parsed for now
// a full parser will be deployed in the future
func ParseConstraint(in string) (Constraint, error) {
	in = strings.TrimSpace(in)
	curr := 0
	l := len(in)
	if l == 0 {
		return &True{}, nil
	}
	next := func() byte {
		if curr < l {
			curr++
			return in[curr-1]
		}
		return 0
	}
	peek := func() byte {
		if curr < l {
			return in[curr]
		}
		return 0
	}

	ver := func() (*Version, error) {
		start := curr
		for {
			n := peek()
			if !isIdentifier(n) && !isVersionSeparator(n) {
				if start == curr {
					return nil, fmt.Errorf("invalid version")
				}
				return Parse(in[start:curr])
			}
			curr++
		}
	}

	stack := []Constraint{}
	for {
		switch next() {
		case '=':
			if v, err := ver(); err == nil {
				stack = append(stack, &Equals{v})
			} else {
				return nil, err
			}
		case '>':
			if peek() == '=' {
				next()
				if v, err := ver(); err == nil {
					stack = append(stack, &GreaterThanOrEqual{v})
				} else {
					return nil, err
				}
			} else {
				if v, err := ver(); err == nil {
					stack = append(stack, &GreaterThan{v})
				} else {
					return nil, err
				}
			}
		case '<':
			if peek() == '=' {
				next()
				if v, err := ver(); err == nil {
					stack = append(stack, &LessThanOrEqual{v})
				} else {
					return nil, err
				}
			} else {
				if v, err := ver(); err == nil {
					stack = append(stack, &LessThan{v})
				} else {
					return nil, err
				}
			}
		case ' ':
			// ignore
		default:
			return nil, fmt.Errorf("unexpected char at: %s", in[curr-1:])
		case 0:
			if len(stack) != 1 {
				return nil, fmt.Errorf("invalid constraint: %s", in)
			}
			return stack[0], nil
		}
	}
}

// True is the empty constraint
type True struct {
}

// Match always return true
func (t *True) Match(v *Version) bool {
	return true
}

func (t *True) String() string {
	return ""
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

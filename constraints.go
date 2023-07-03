//
// Copyright 2018-2023 Cristian Maglie. All rights reserved.
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
	skipSpace := func() {
		for curr < l && in[curr] == ' ' {
			curr++
		}
	}
	peek := func() byte {
		if curr < l {
			return in[curr]
		}
		return 0
	}

	version := func() (*Version, error) {
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

	var terminal func() (Constraint, error)
	var constraint func() (Constraint, error)

	terminal = func() (Constraint, error) {
		skipSpace()
		switch next() {
		case '!':
			expr, err := terminal()
			if err != nil {
				return nil, err
			}
			return &Not{expr}, nil
		case '(':
			expr, err := constraint()
			if err != nil {
				return nil, err
			}
			skipSpace()
			if c := next(); c != ')' {
				return nil, fmt.Errorf("unexpected char at: %s", in[curr-1:])
			}
			return expr, nil
		case '=':
			v, err := version()
			if err != nil {
				return nil, err
			}
			return &Equals{v}, nil
		case '^':
			v, err := version()
			if err != nil {
				return nil, err
			}
			return &CompatibleWith{v}, nil
		case '>':
			if peek() == '=' {
				next()
				v, err := version()
				if err != nil {
					return nil, err
				}
				return &GreaterThanOrEqual{v}, nil
			} else {
				v, err := version()
				if err != nil {
					return nil, err
				}
				return &GreaterThan{v}, nil
			}
		case '<':
			if peek() == '=' {
				next()
				v, err := version()
				if err != nil {
					return nil, err
				}
				return &LessThanOrEqual{v}, nil
			} else {
				v, err := version()
				if err != nil {
					return nil, err
				}
				return &LessThan{v}, nil
			}
		default:
			return nil, fmt.Errorf("unexpected char at: %s", in[curr-1:])
		}
	}

	andExpr := func() (Constraint, error) {
		t1, err := terminal()
		if err != nil {
			return nil, err
		}
		stack := []Constraint{t1}

		for {
			skipSpace()
			if peek() != '&' {
				if len(stack) == 1 {
					return stack[0], nil
				}
				return &And{stack}, nil
			}
			next()
			if peek() != '&' {
				return nil, fmt.Errorf("unexpected char at: %s", in[curr-1:])
			}
			next()

			t2, err := terminal()
			if err != nil {
				return nil, err
			}
			stack = append(stack, t2)
		}
	}

	constraint = func() (Constraint, error) {
		t1, err := andExpr()
		if err != nil {
			return nil, err
		}
		stack := []Constraint{t1}

		for {
			skipSpace()
			switch peek() {
			case '|':
				next()
				if peek() != '|' {
					return nil, fmt.Errorf("unexpected char at: %s", in[curr-1:])
				}
				next()

				t2, err := andExpr()
				if err != nil {
					return nil, err
				}
				stack = append(stack, t2)

			case 0, ')':
				if len(stack) == 1 {
					return stack[0], nil
				}
				return &Or{stack}, nil

			default:
				return nil, fmt.Errorf("unexpected char at: %s", in[curr-1:])
			}
		}
	}

	return constraint()
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

// CompatibleWith is the "compatible with" (^) constraint
type CompatibleWith struct {
	Version *Version
}

// Match returns true if v satisfies the condition
func (cw *CompatibleWith) Match(v *Version) bool {
	return cw.Version.CompatibleWith(v)
}

func (cw *CompatibleWith) String() string {
	return "^" + cw.Version.String()
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

// Not match if Operand does not match and viceversa
type Not struct {
	Operand Constraint
}

// Match returns ture if v does NOT satisfies the condition
func (not *Not) Match(v *Version) bool {
	return !not.Operand.Match(v)
}

func (not *Not) String() string {
	op := not.Operand.String()
	if op == "" || op[0] != '(' {
		return "!(" + op + ")"
	}
	return "!" + op
}

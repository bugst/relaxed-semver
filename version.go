//
// Copyright 2018-2025 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

// Version contains the results of parsed version string
type Version struct {
	raw        string
	bytes      []byte
	major      int
	minor      int
	patch      int
	prerelease int
	build      int
}

func (v *Version) String() string {
	if v == nil {
		return ""
	}
	return v.raw
}

// NormalizedString is a datatype to be used in maps and other places where the
// version is used as a key.
type NormalizedString string

// NormalizedString return a string representation of the version that is
// normalized to always have a major, minor and patch version. This is useful
// to be used in maps and other places where the version is used as a key.
func (v *Version) NormalizedString() NormalizedString {
	if v == nil {
		return ""
	}
	if v.major == 0 {
		return NormalizedString("0.0.0")
	} else if v.minor == v.major {
		return NormalizedString(v.raw[0:v.major] + ".0.0" + v.raw[v.major:])
	} else if v.patch == v.minor {
		return NormalizedString(v.raw[0:v.minor] + ".0" + v.raw[v.minor:])
	} else {
		return NormalizedString(v.raw)
	}
}

// Normalize transforms a truncated semver version in a strictly compliant semver
// version by adding minor and patch versions. For example:
// "1" is trasformed to "1.0.0" or "2.5-dev" to "2.5.0-dev"
func (v *Version) Normalize() {
	if v.major == 0 {
		v.raw = "0.0.0" + v.raw
		v.major = 1
		v.minor = 3
		v.patch = 5
		v.prerelease += 5
		v.build += 5
	} else if v.minor == v.major {
		v.raw = v.raw[0:v.major] + ".0.0" + v.raw[v.major:]
		v.minor = v.major + 2
		v.patch = v.major + 4
		v.prerelease += 4
		v.build += 4
	} else if v.patch == v.minor {
		v.raw = v.raw[0:v.minor] + ".0" + v.raw[v.minor:]
		v.patch = v.minor + 2
		v.prerelease += 2
		v.build += 2
	}
	v.bytes = []byte(v.raw)
}

func compareNumber(a, b []byte) int {
	la := len(a)
	lb := len(b)
	if la == lb {
		for i := range a {
			if a[i] == b[i] {
				continue
			}
			if a[i] > b[i] {
				return 1
			}
			return -1
		}
		return 0
	}
	if la > lb {
		return 1
	}
	return -1
}

func compareAlpha(a, b []byte) int {
	if string(a) > string(b) {
		return 1
	}
	if string(a) < string(b) {
		return -1
	}
	return 0
}

var zero = []byte("0")

// CompareTo compares the Version with the one passed as parameter.
// Returns -1, 0 or 1 if the version is respectively less than, equal
// or greater than the compared Version
func (v *Version) CompareTo(u *Version) int {
	// 11. Precedence refers to how versions are compared to each other when ordered.
	// Precedence MUST be calculated by separating the version into cmp, minor,
	// patch and pre-release identifiers in that order (Build metadata does not
	// figure into precedence). Precedence is determined by the first difference when
	// comparing each of these identifiers from left to right as follows: Major, minor,
	// and patch versions are always compared numerically.
	// Example: 1.0.0 < 2.0.0 < 2.1.0 < 2.1.1.
	vIdx := 0
	uIdx := 0
	vMajor := v.major
	uMajor := u.major
	{
		if vMajor == uMajor {
			for vIdx < vMajor {
				if v.bytes[vIdx] == u.bytes[uIdx] {
					vIdx++
					uIdx++
					continue
				}
				if v.bytes[vIdx] > u.bytes[uIdx] {
					return 1
				}
				return -1
			}
		} else if vMajor == 0 && u.bytes[uIdx] == '0' {
			// continue
		} else if uMajor == 0 && v.bytes[vIdx] == '0' {
			// continue
		} else if vMajor > uMajor {
			return 1
		} else {
			return -1
		}
	}
	vMinor := v.minor
	uMinor := u.minor
	vIdx = vMajor + 1
	uIdx = uMajor + 1
	{
		la := vMinor - vMajor - 1
		lb := uMinor - uMajor - 1
		if la == lb {
			for vIdx < vMinor {
				if v.bytes[vIdx] == u.bytes[uIdx] {
					vIdx++
					uIdx++
					continue
				}
				if v.bytes[vIdx] > u.bytes[uIdx] {
					return 1
				}
				return -1
			}
		} else if vMinor == vMajor && u.bytes[uIdx] == '0' {
			// continue
		} else if uMinor == uMajor && v.bytes[vIdx] == '0' {
			// continue
		} else if la > lb {
			return 1
		} else {
			return -1
		}
	}
	vPatch := v.patch
	uPatch := u.patch
	vIdx = vMinor + 1
	uIdx = uMinor + 1
	{
		la := vPatch - vMinor - 1
		lb := uPatch - uMinor - 1
		if la == lb {
			for vIdx < vPatch {
				if v.bytes[vIdx] == u.bytes[uIdx] {
					vIdx++
					uIdx++
					continue
				}
				if v.bytes[vIdx] > u.bytes[uIdx] {
					return 1
				}
				return -1
			}
		} else if vPatch == vMinor && u.bytes[uIdx] == '0' {
			// continue
		} else if uPatch == uMinor && v.bytes[vIdx] == '0' {
			// continue
		} else if la > lb {
			return 1
		} else {
			return -1
		}
	}

	// if both versions have no pre-release, they are equal
	if v.prerelease == vPatch && u.prerelease == uPatch {
		return 0
	}

	// When major, minor, and patch are equal, a pre-release version has lower
	// precedence than a normal version.
	// Example: 1.0.0-alpha < 1.0.0.

	// if v has no pre-release, it's greater than u
	if v.prerelease == vPatch {
		return 1
	}
	// if u has no pre-release, it's greater than v
	if u.prerelease == uPatch {
		return -1
	}

	// Precedence for two pre-release versions with the same major, minor, and patch
	// version MUST be determined by comparing each dot separated identifier from left
	// to right until a difference is found as follows:
	// - identifiers consisting of only digits are compared numerically
	// - identifiers with letters or hyphens are compared lexically in ASCII sort order.
	// Numeric identifiers always have lower precedence than non-numeric identifiers.
	// A larger set of pre-release fields has a higher precedence than a smaller set,
	// if all of the preceding identifiers are equal.
	// Example: 1.0.0-alpha < 1.0.0-alpha.1 < 1.0.0-alpha.beta < 1.0.0-beta <
	//          < 1.0.0-beta.2 < 1.0.0-beta.11 < 1.0.0-rc.1 < 1.0.0.
	vIdx = vPatch + 1
	uIdx = uPatch + 1
	vLast := v.prerelease
	uLast := u.prerelease
	vIsAlpha := false
	uIsAlpha := false
	vIsLonger := false
	uIsLonger := false
	cmp := 0
	for {
		var vCurr byte
		var uCurr byte
		if vIdx != vLast {
			vCurr = v.raw[vIdx]
		}
		if uIdx != uLast {
			uCurr = u.raw[uIdx]
		}

		if vIdx == vLast || vCurr == '.' {
			if uIdx != uLast && uCurr != '.' {
				if !uIsAlpha && !(uCurr >= '0' && uCurr <= '9') {
					uIsAlpha = true
				}
				uIsLonger = true
				uIdx++
				continue
			}
		} else if uIdx == uLast || uCurr == '.' {
			if vIdx != vLast && vCurr != '.' {
				if !vIsAlpha && !(vCurr >= '0' && vCurr <= '9') {
					vIsAlpha = true
				}
				vIsLonger = true
				vIdx++
				continue
			}
		} else {
			if cmp == 0 {
				if vCurr > uCurr {
					cmp = 1
				} else if vCurr < uCurr {
					cmp = -1
				}
			}
			if !vIsAlpha && !(vCurr >= '0' && vCurr <= '9') {
				vIsAlpha = true
			}
			if !uIsAlpha && !(uCurr >= '0' && uCurr <= '9') {
				uIsAlpha = true
			}
			vIdx++
			uIdx++
			continue
		}

		// Numeric identifiers always have lower precedence than non-numeric identifiers.
		if vIsAlpha && uIsAlpha {
			if cmp != 0 {
				// alphanumeric vs alphanumeric, sorting has priority
				return cmp
			} else if vIsLonger {
				// alphanumeric vs alphanumeric, v is longer, return >
				return 1
			} else if uIsLonger {
				// alphanumeric vs alphanumeric, u is longer, return <
				return -1
			}
			// Both alphanumeric, if comparison is equal, move on the next field
		} else if vIsAlpha && !uIsAlpha {
			// alphanumeric vs numeric, return >
			return 1
		} else if !vIsAlpha && uIsAlpha {
			// numeric vs alphanumeric, return <
			return -1
		} else {
			if vIsLonger {
				// numeric vs numeric, v is longer, return >
				return 1
			} else if uIsLonger {
				// numeric vs numeric, u is longer, return <
				return -1
			} else if cmp != 0 {
				// numeric vs numeric, return cmp if not equal
				return cmp
			}
			// Both numeric, if comparison is equal, move on the next field
		}

		// A larger set of pre-release fields has a higher precedence than a smaller set,
		// if all of the preceding identifiers are equal.

		if vIdx == vLast && uIdx == uLast {
			// No more field, proceed with build metadata
			break
		}
		if vIdx != vLast && uIdx == uLast {
			// v has more fields, return >
			return 1
		}
		if vIdx == vLast && uIdx != uLast {
			// u has more fields, return <
			return -1
		}

		// Move on the next field
		vIsAlpha = false
		uIsAlpha = false
		vIsLonger = false
		uIsLonger = false
		vIdx++
		uIdx++
	}

	return 0
}

// LessThan returns true if the Version is less than the Version passed as parameter
func (v *Version) LessThan(u *Version) bool {
	return v.CompareTo(u) < 0
}

// LessThanOrEqual returns true if the Version is less than or equal to the Version passed as parameter
func (v *Version) LessThanOrEqual(u *Version) bool {
	return v.CompareTo(u) <= 0
}

// Equal returns true if the Version is equal to the Version passed as parameter
func (v *Version) Equal(u *Version) bool {
	return v.CompareTo(u) == 0
}

// GreaterThan returns true if the Version is greater than the Version passed as parameter
func (v *Version) GreaterThan(u *Version) bool {
	return v.CompareTo(u) > 0
}

// GreaterThanOrEqual returns true if the Version is greater than or equal to the Version passed as parameter
func (v *Version) GreaterThanOrEqual(u *Version) bool {
	return v.CompareTo(u) >= 0
}

// CompatibleWith returns true if the Version is compatible with the version passed as paramater
func (v *Version) CompatibleWith(u *Version) bool {
	if !u.GreaterThanOrEqual(v) {
		return false
	}
	vMajor := zero[:]
	if v.major > 0 {
		vMajor = v.bytes[:v.major]
	}
	uMajor := zero[:]
	if u.major > 0 {
		uMajor = u.bytes[:u.major]
	}
	majorEquals := compareNumber(vMajor, uMajor) == 0
	if v.major > 0 && v.bytes[0] != '0' {
		return majorEquals
	}
	if !majorEquals {
		return false
	}
	vMinor := zero[:]
	if v.minor > v.major {
		vMinor = v.bytes[v.major+1 : v.minor]
	}
	uMinor := zero[:]
	if u.minor > u.major {
		uMinor = u.bytes[u.major+1 : u.minor]
	}
	minorEquals := compareNumber(vMinor, uMinor) == 0
	if vMinor[0] != '0' {
		return minorEquals
	}
	if !minorEquals {
		return false
	}
	vPatch := zero[:]
	if v.patch > v.minor {
		vPatch = v.bytes[v.minor+1 : v.patch]
	}
	uPatch := zero[:]
	if u.patch > u.minor {
		uPatch = u.bytes[u.minor+1 : u.patch]
	}
	return compareNumber(vPatch, uPatch) == 0
}

// SortableString returns the version encoded as a string that when compared
// with alphanumeric ordering it respects the original semver ordering:
//
//	(v1.SortableString() < v2.SortableString()) == v1.LessThan(v2)
//	cmp.Compare[string](v1.SortableString(), v2.SortableString()) == v1.CompareTo(v2)
//
// This may turn out useful when the version is saved in a database or is
// introduced in a system that doesn't support semver ordering.
func (v *Version) SortableString() string {
	// Encode a number in a string that when compared as string it respects
	// the original numeric order.
	// To allow longer numbers to be compared correctly, a prefix of ":"s
	// with the length of the number is added minus 1.
	// For example: 123 -> "::123"
	//              45  -> ":45"
	// The number written as string compare as ("123" < "99") but the encoded
	// version keeps the original integer ordering ("::123" > ":99").
	encodeNumber := func(in []byte) string {
		if len(in) == 0 {
			return "0"
		}
		p := ""
		for range in {
			p += ":"
		}
		return p[:len(p)-1] + string(in)
	}

	var vMajor, vMinor, vPatch []byte
	vMajor = v.bytes[:v.major]
	if v.minor > v.major {
		vMinor = v.bytes[v.major+1 : v.minor]
	}
	if v.patch > v.minor {
		vPatch = v.bytes[v.minor+1 : v.patch]
	}

	res := encodeNumber(vMajor) + "." + encodeNumber(vMinor) + "." + encodeNumber(vPatch)
	// If there is no pre-release, add a ";" to the end, otherwise add a "-" followed by the pre-release.
	// This ensure the correct ordering of the pre-release versions (that are always lower than the normal versions).
	if v.prerelease == v.patch {
		return res + ";"
	}
	res += "-"

	isAlpha := false
	add := func(in []byte) {
		// if the pre-release piece is alphanumeric, add a ";" before the piece
		// otherwise add an ":" before the piece. This ensure the correct ordering
		// of the pre-release piece (numeric are lower than alphanumeric).
		if isAlpha {
			res += ";" + string(in)
		} else {
			res += ":" + encodeNumber(in)
		}
		isAlpha = false
	}
	prerelease := v.bytes[v.patch+1 : v.prerelease]
	start := 0
	for curr, c := range prerelease {
		if c == '.' {
			add(prerelease[start:curr])
			// separate the pre-release pieces with a "," to ensure the correct ordering
			// of the pre-release pieces (the separator must be lower than any other allowed
			// character [a-zA-Z0-9-]).
			res += ","
			start = curr + 1
			continue
		}
		if !isNumeric(c) {
			isAlpha = true
		}
	}
	add(prerelease[start:])
	return res
}

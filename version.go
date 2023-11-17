//
// Copyright 2018-2023 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

// Version contains the results of parsed version string
type Version struct {
	raw        string
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
}

func compareNumberRelaxed(a, b string) int {
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

	// la != lb...

	// in relaxed semver, a missing number is considered 0
	if la == 0 {
		if b[0] == '0' {
			return 0 // ("","0")
		}
		return -1 // ("","N")
	}
	if lb == 0 {
		if a[0] == '0' {
			return 0 // ("0","")
		}
		return 1 // ("N","")
	}

	// both a and b are not empty
	if la > lb {
		return 1 // ("NNN","N")
	}
	return -1 // ("N","NNN")
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
	if cmp := compareNumberRelaxed(v.majorString(), u.majorString()); cmp != 0 {
		return cmp
	}
	if cmp := compareNumberRelaxed(v.minorString(), u.minorString()); cmp != 0 {
		return cmp
	}
	if cmp := compareNumberRelaxed(v.patchString(), u.patchString()); cmp != 0 {
		return cmp
	}

	// if both versions have no pre-release, they are equal
	if v.prerelease == v.patch && u.prerelease == u.patch {
		return 0
	}

	// When major, minor, and patch are equal, a pre-release version has lower
	// precedence than a normal version.
	// Example: 1.0.0-alpha < 1.0.0.

	// if v has no pre-release, it's greater than u
	if v.prerelease == v.patch {
		return 1
	}
	// if u has no pre-release, it's greater than v
	if u.prerelease == u.patch {
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
	vIdx := v.patch + 1
	uIdx := u.patch + 1
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
				if !uIsAlpha && !isNumeric(uCurr) {
					uIsAlpha = true
				}
				uIsLonger = true
				uIdx++
				continue
			}
		} else if uIdx == uLast || uCurr == '.' {
			if vIdx != vLast && vCurr != '.' {
				if !vIsAlpha && !isNumeric(vCurr) {
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
			if !vIsAlpha && !isNumeric(vCurr) {
				vIsAlpha = true
			}
			if !uIsAlpha && !isNumeric(uCurr) {
				uIsAlpha = true
			}
			vIdx++
			uIdx++
			continue
		}

		// Numeric identifiers always have lower precedence than non-numeric identifiers.
		if vIsAlpha && uIsAlpha {
			if cmp != 0 {
				return cmp
			}
			// Both alphanumeric, if comparison is equal, move on the next field
		} else if vIsAlpha && !uIsAlpha {
			// alphanumeric vs numeric, return >
			return 1
		} else if !vIsAlpha && uIsAlpha {
			// numeric vs alphanumeric, return <
			return -1
		} else if vIsLonger {
			// numeric vs numeric, v is longer, return >
			return 1
		} else if uIsLonger {
			// numeric vs numeric, u is longer, return <
			return -1
		} else if cmp != 0 {
			// numeric vs numeric, return cmp if not equal
			return cmp
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
	majorEquals := compareNumberRelaxed(v.majorString(), u.majorString()) == 0
	if v.major > 0 && v.raw[0] != '0' {
		return majorEquals
	}
	if !majorEquals {
		return false
	}
	minorEquals := compareNumberRelaxed(v.minorString(), u.minorString()) == 0
	if v.minor > v.major && v.raw[v.major+1] != '0' {
		return minorEquals
	}
	if !minorEquals {
		return false
	}
	return compareNumberRelaxed(v.patchString(), u.patchString()) == 0
}

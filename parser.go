//
// Copyright 2018-2023 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"fmt"
)

var empty = []byte("")

// MustParse parse a version string and panic if the parsing fails
func MustParse(inVersion string) *Version {
	res, err := Parse(inVersion)
	if err != nil {
		panic(err)
	}
	return res
}

// Parse parse a version string
func Parse(inVersion string) (*Version, error) {
	result := &Version{
		major: empty[:],
		minor: empty[:],
		patch: empty[:],
	}
	if err := parseInto([]byte(inVersion), result); err != nil {
		return nil, err
	}
	return result, nil
}

func parseInto(in []byte, result *Version) error {
	// Setup parsing harness
	inLen := len(in)
	currIdx := -1
	var curr byte
	next := func() bool {
		currIdx = currIdx + 1
		if currIdx == inLen {
			return false
		}
		curr = in[currIdx]
		return true
	}

	// 2. A normal version number MUST take the form X.Y.Z where X, Y, and Z
	// are non-negative integers, and MUST NOT contain leading zeroes. X is
	// the major version, Y is the minor version, and Z is the patch version.
	// Each element MUST increase numerically.
	// For instance: 1.9.0 -> 1.10.0 -> 1.11.0.

	// Parse major
	if !next() {
		return nil // empty version
	}
	if !numeric[curr] {
		return fmt.Errorf("no major version found")
	}
	if curr == '0' {
		result.major = in[0:1] // 0
		if !next() {
			return nil
		}
		if numeric[curr] {
			return fmt.Errorf("major version must not be prefixed with zero")
		}
		if !versionSeparator[curr] {
			return fmt.Errorf("invalid major version separator '%c'", curr)
		}
		// Fallthrough and parse next element
	} else {
		for {
			if !next() {
				result.major = in[0:currIdx]
				return nil
			}
			if numeric[curr] {
				continue
			}
			if versionSeparator[curr] {
				result.major = in[0:currIdx]
				break
			}
			return fmt.Errorf("invalid major version separator '%c'", curr)
		}
	}

	// Parse minor
	if curr == '.' {
		if !next() || !numeric[curr] {
			return fmt.Errorf("no minor version found")
		}
		if curr == '0' {
			result.minor = in[currIdx : currIdx+1] // x.0
			if !next() {
				return nil
			}
			if numeric[curr] {
				return fmt.Errorf("minor version must not be prefixed with zero")
			}
			if !versionSeparator[curr] {
				return fmt.Errorf("invalid minor version separator '%c'", curr)
			}
			// Fallthrough and parse next element
		} else {
			minorIdx := currIdx
			for {
				if !next() {
					result.minor = in[minorIdx:currIdx]
					return nil
				}
				if numeric[curr] {
					continue
				}
				if versionSeparator[curr] {
					result.minor = in[minorIdx:currIdx]
					break
				}
				return fmt.Errorf("invalid minor version separator '%c'", curr)
			}
		}
	}

	// Parse patch
	if curr == '.' {
		if !next() || !numeric[curr] {
			return fmt.Errorf("no patch version found")
		}
		if curr == '0' {
			result.patch = in[currIdx : currIdx+1] // x.y.0
			if !next() {
				return nil
			}
			if numeric[curr] {
				return fmt.Errorf("patch version must not be prefixed with zero")
			}
			if !versionSeparator[curr] {
				return fmt.Errorf("invalid patch version separator '%c'", curr)
			}
			// Fallthrough and parse next element
		} else {
			patchIdx := currIdx
			for {
				if !next() {
					result.patch = in[patchIdx:currIdx]
					return nil
				}
				if numeric[curr] {
					continue
				}
				if curr == '-' || curr == '+' {
					result.patch = in[patchIdx:currIdx]
					break
				}
				return fmt.Errorf("invalid patch version separator '%c'", curr)
			}
		}
	}

	// 9. A pre-release version MAY be denoted by appending a hyphen and a series
	// of dot separated identifiers immediately following the patch version.
	// Identifiers MUST comprise only ASCII alphanumerics and hyphen [0-9A-Za-z-].
	// Identifiers MUST NOT be empty.
	// Numeric identifiers MUST NOT include leading zeroes.
	// Pre-release versions have a lower precedence than the associated normal
	// version. A pre-release version indicates that the version is unstable and
	// might not satisfy the intended compatibility requirements as denoted by
	// its associated normal version.
	// Examples: 1.0.0-alpha, 1.0.0-alpha.1, 1.0.0-0.3.7, 1.0.0-x.7.z.92.
	if curr == '-' {
		// Pre-release parsing

		prereleaseIdx := currIdx + 1
		zeroPrefix := false
		alphaIdentifier := false
		for {
			if hasNext := next(); !hasNext || curr == '.' || curr == '+' {
				if prereleaseIdx == currIdx {
					return fmt.Errorf("empty prerelease not allowed")
				}
				if zeroPrefix && !alphaIdentifier && currIdx-prereleaseIdx > 1 {
					return fmt.Errorf("numeric prerelease must not be prefixed with zero")
				}
				result.prerelases = append(result.prerelases, in[prereleaseIdx:currIdx])
				result.numericPrereleases = append(result.numericPrereleases, !alphaIdentifier)
				if !hasNext {
					return nil
				}
				if curr == '+' {
					break
				}

				// Multiple prerelease
				prereleaseIdx = currIdx + 1
				zeroPrefix = false
				alphaIdentifier = false
				continue
			}
			if prereleaseIdx == currIdx && curr == '0' {
				zeroPrefix = true
				continue
			}
			if numeric[curr] {
				continue
			}
			if identifier[curr] {
				alphaIdentifier = true
				continue
			}
			return fmt.Errorf("invalid prerelease separator: '%c'", curr)
		}
	}

	// 10. Build metadata MAY be denoted by appending a plus sign and a series of
	// dot separated identifiers immediately following the patch or pre-release
	// version.
	// Identifiers MUST comprise only ASCII alphanumerics and hyphen [0-9A-Za-z-].
	// Identifiers MUST NOT be empty.
	// Build metadata SHOULD be ignored when determining version precedence. Thus
	// two versions that differ only in the build metadata, have the same precedence.
	// Examples: 1.0.0-alpha+001, 1.0.0+20130313144700, 1.0.0-beta+exp.sha.5114f85.

	// Builds parsing
	buildIdx := currIdx + 1
	if curr == '+' {
		for {
			if hasNext := next(); !hasNext || curr == '.' {
				if buildIdx == currIdx {
					return fmt.Errorf("empty build tag not allowed")
				}
				result.builds = append(result.builds, in[buildIdx:currIdx])
				if !hasNext {
					return nil
				}

				// Multiple builds
				buildIdx = currIdx + 1
				continue
			}
			if identifier[curr] {
				continue
			}
			return fmt.Errorf("invalid separator for builds: '%c'", curr)
		}
	}
	return fmt.Errorf("invalid separator: '%c'", curr)
}

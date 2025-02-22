//
// Copyright 2018-2025 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"gopkg.in/yaml.v3"
)

// MarshalYAML implements yaml.Marshaler
func (v *Version) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaler
func (v *Version) UnmarshalYAML(node *yaml.Node) error {
	var versionString string
	if err := node.Decode(&versionString); err != nil {
		return err
	}
	parsed, err := Parse(versionString)
	if err != nil {
		return err
	}

	v.raw = parsed.raw
	v.bytes = []byte(v.raw)
	v.major = parsed.major
	v.minor = parsed.minor
	v.patch = parsed.patch
	v.prerelease = parsed.prerelease
	v.build = parsed.build
	return nil
}

// MarshalYAML implements yaml.Marshaler
func (v *RelaxedVersion) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaler
func (v *RelaxedVersion) UnmarshalYAML(node *yaml.Node) error {
	var versionString string
	if err := node.Decode(&versionString); err != nil {
		return err
	}
	parsed := ParseRelaxed(versionString)

	v.customversion = parsed.customversion
	v.version = parsed.version
	return nil
}

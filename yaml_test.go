//
// Copyright 2018-2023 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestYAMLParseVersion(t *testing.T) {
	var versionIsYamlUnmarshaler yaml.Unmarshaler = MustParse("1.0.0")
	var versionIsYamlMarshaler yaml.Marshaler = MustParse("1.0.0")
	_ = versionIsYamlUnmarshaler
	_ = versionIsYamlMarshaler

	testVersion := "1.2.3-aaa.4.5.6+bbb.7.8.9"
	v, err := Parse(testVersion)
	require.NoError(t, err)

	data, err := yaml.Marshal(v)
	require.Equal(t, "1.2.3-aaa.4.5.6+bbb.7.8.9\n", string(data))
	require.NoError(t, err)

	var u Version
	err = yaml.Unmarshal(data, &u)
	require.NoError(t, err)

	dump := fmt.Sprintf("%v,%v,%v,%v,%v,%v", u.raw, u.major, u.minor, u.patch, u.prerelease, u.build)
	require.Equal(t, "1.2.3-aaa.4.5.6+bbb.7.8.9,1,3,5,15,25", dump)

	require.Equal(t, testVersion, u.String())

	err = yaml.Unmarshal([]byte(`"invalid"`), &u)
	require.Error(t, err)

	err = yaml.Unmarshal([]byte(`invalid:`), &u)
	require.Error(t, err)

	require.NoError(t, yaml.Unmarshal([]byte(`"1.6.2"`), &v))
	require.NoError(t, yaml.Unmarshal([]byte(`"1.6.3"`), &u))
	require.True(t, u.GreaterThan(v))
}

func TestYAMLParseRelaxedVersion(t *testing.T) {
	var relaxedVersionIsYamlUnmarshaler yaml.Unmarshaler = ParseRelaxed("1.0.0")
	var relaxedVersionIsYamlMarshaler yaml.Marshaler = ParseRelaxed("1.0.0")
	_ = relaxedVersionIsYamlUnmarshaler
	_ = relaxedVersionIsYamlMarshaler

	testVersion := "1.2.3-aaa.4.5.6+bbb.7.8.9"

	v := ParseRelaxed(testVersion)

	data, err := yaml.Marshal(v)
	require.NoError(t, err)
	require.Equal(t, "1.2.3-aaa.4.5.6+bbb.7.8.9\n", string(data))

	var u RelaxedVersion
	err = yaml.Unmarshal(data, &u)
	require.NoError(t, err)

	require.Equal(t, testVersion, u.String())

	err = yaml.Unmarshal([]byte(`"invalid"`), &u)
	require.NoError(t, err)
	require.Equal(t, "invalid", u.String())

	err = yaml.Unmarshal([]byte(`invalid:`), &u)
	require.Error(t, err)
}

func BenchmarkYAMLDecoding(b *testing.B) {
	testVersion := "1.2.3-aaa.4.5.6+bbb.7.8.9"
	v, _ := Parse(testVersion)
	data, _ := yaml.Marshal(v)
	var u Version
	for i := 0; i < b.N; i++ {
		yaml.Unmarshal(data, &u)
	}
}

func BenchmarkYAMLDecodingRelaxed(b *testing.B) {
	testVersion := "1.2.3-aaa.4.5.6+bbb.7.8.9"
	v := ParseRelaxed(testVersion)
	data, _ := yaml.Marshal(v)
	var u RelaxedVersion
	for i := 0; i < b.N; i++ {
		yaml.Unmarshal(data, &u)
	}
}

//
// Copyright 2018-2022 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSONParseVersion(t *testing.T) {
	testVersion := "1.2.3-aaa.4.5.6+bbb.7.8.9"
	v, err := Parse(testVersion)
	require.NoError(t, err)

	data, err := json.Marshal(v)
	fmt.Println(string(data))
	require.NoError(t, err)

	dump := fmt.Sprintf("%s,%s,%s,%s,%v,%s",
		v.major, v.minor, v.patch,
		v.prerelases, v.numericPrereleases,
		v.builds)
	require.Equal(t, "1,2,3,[aaa 4 5 6],[false true true true],[bbb 7 8 9]", dump)

	var u Version
	err = json.Unmarshal(data, &u)
	require.NoError(t, err)

	require.Equal(t, testVersion, v.String())

	err = json.Unmarshal([]byte(`"invalid"`), &u)
	require.Error(t, err)

	err = json.Unmarshal([]byte(`123`), &u)
	require.Error(t, err)
}

func TestJSONParseRelaxedVersion(t *testing.T) {
	testVersion := "1.2.3-aaa.4.5.6+bbb.7.8.9"
	v := ParseRelaxed(testVersion)

	data, err := json.Marshal(v)
	fmt.Println(string(data))
	require.NoError(t, err)

	var u RelaxedVersion
	err = json.Unmarshal(data, &u)
	require.NoError(t, err)

	require.Equal(t, testVersion, v.String())

	err = json.Unmarshal([]byte(`"invalid"`), &u)
	require.NoError(t, err)
	require.Equal(t, "invalid", u.String())

	err = json.Unmarshal([]byte(`123`), &u)
	require.Error(t, err)
}

func BenchmarkJSONDecoding(b *testing.B) {
	testVersion := "1.2.3-aaa.4.5.6+bbb.7.8.9"
	v, _ := Parse(testVersion)
	data, _ := json.Marshal(v)
	var u Version
	for i := 0; i < b.N; i++ {
		json.Unmarshal(data, &u)
	}
}

func BenchmarkJSONDecodingRelaxed(b *testing.B) {
	testVersion := "1.2.3-aaa.4.5.6+bbb.7.8.9"
	v := ParseRelaxed(testVersion)
	data, _ := json.Marshal(v)
	var u RelaxedVersion
	for i := 0; i < b.N; i++ {
		json.Unmarshal(data, &u)
	}
}

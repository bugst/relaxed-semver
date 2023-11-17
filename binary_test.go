//
// Copyright 2018-2023 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGOBEncoderVersion(t *testing.T) {
	testVersion := "1.2.3-aaa.4.5.6+bbb.7.8.9"

	v, err := Parse(testVersion)
	require.NoError(t, err)
	dumpV := fmt.Sprintf("%v,%v,%v,%v,%v,%v", v.raw, v.major, v.minor, v.patch, v.prerelease, v.build)
	require.Equal(t, "1.2.3-aaa.4.5.6+bbb.7.8.9,1,3,5,15,25", dumpV)
	require.Equal(t, testVersion, v.String())

	dataV := new(bytes.Buffer)
	err = gob.NewEncoder(dataV).Encode(v)
	require.NoError(t, err)

	var u Version
	err = gob.NewDecoder(dataV).Decode(&u)
	require.NoError(t, err)
	dumpU := fmt.Sprintf("%v,%v,%v,%v,%v,%v", v.raw, u.major, u.minor, u.patch, u.prerelease, u.build)

	require.Equal(t, dumpV, dumpU)
	require.Equal(t, testVersion, u.String())
}

func TestGOBEncoderRelaxedVersion(t *testing.T) {
	check := func(testVersion string) {
		v := ParseRelaxed(testVersion)

		dataV := new(bytes.Buffer)
		err := gob.NewEncoder(dataV).Encode(v)
		require.NoError(t, err)

		var u RelaxedVersion
		err = gob.NewDecoder(dataV).Decode(&u)
		require.NoError(t, err)

		require.Equal(t, testVersion, u.String())
	}
	check("1.2.3-aaa.4.5.6+bbb.7.8.9")
	check("asdasdasd-1.2.3-aaa.4.5.6+bbb.7.8.9")
}

func BenchmarkBinaryDecoding(b *testing.B) {
	testVersion := "1.2.3-aaa.4.5.6+bbb.7.8.9"
	v := MustParse(testVersion)

	data, _ := v.MarshalBinary()
	var u Version
	for i := 0; i < b.N; i++ {
		u.UnmarshalBinary(data)
	}
}

func BenchmarkBinaryDecodingRelaxed(b *testing.B) {
	testVersion := "1.2.3-aaa.4.5.6+bbb.7.8.9"
	v := ParseRelaxed(testVersion)

	data, _ := v.MarshalBinary()
	var u RelaxedVersion
	for i := 0; i < b.N; i++ {
		u.UnmarshalBinary(data)
	}
}

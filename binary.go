//
// Copyright 2018-2023 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"bytes"
	"encoding/binary"
)

func marshalByteArray(b []byte) []byte {
	l := len(b)
	res := make([]byte, l+4)
	binary.BigEndian.PutUint32(res, uint32(l))
	copy(res[4:], b)
	return res
}

func marshalInt(i int) []byte {
	res := make([]byte, 4)
	binary.BigEndian.PutUint32(res, uint32(i))
	return res
}

// MarshalBinary implements binary custom encoding
func (v *Version) MarshalBinary() ([]byte, error) {
	res := new(bytes.Buffer)
	_, _ = res.Write(marshalByteArray(v.major))
	_, _ = res.Write(marshalByteArray(v.minor))
	_, _ = res.Write(marshalByteArray(v.patch))
	_, _ = res.Write(marshalInt(len(v.prerelases)))
	for _, pre := range v.prerelases {
		_, _ = res.Write(marshalByteArray(pre))
	}
	_, _ = res.Write(marshalInt(len(v.numericPrereleases)))
	for _, npre := range v.numericPrereleases {
		v := []byte{0}
		if npre {
			v[0] = 1
		}
		_, _ = res.Write(v)
	}
	_, _ = res.Write(marshalInt(len(v.builds)))
	for _, build := range v.builds {
		_, _ = res.Write(marshalByteArray(build))
	}
	return res.Bytes(), nil
}

func decodeArray(data []byte) ([]byte, []byte) {
	l, data := int(binary.BigEndian.Uint32(data)), data[4:]
	return data[:l], data[l:]
}

func decodeInt(data []byte) (int, []byte) {
	return int(binary.BigEndian.Uint32(data)), data[4:]
}

// UnmarshalJSON implements binary custom decoding
func (v *Version) UnmarshalBinary(data []byte) error {
	var buff []byte

	v.major, data = decodeArray(data)
	v.minor, data = decodeArray(data)
	v.patch, data = decodeArray(data)
	n, data := decodeInt(data)
	v.prerelases = nil
	for i := 0; i < n; i++ {
		buff, data = decodeArray(data)
		v.prerelases = append(v.prerelases, buff)
	}
	v.numericPrereleases = nil
	n, data = decodeInt(data)
	for i := 0; i < n; i++ {
		num := false
		if data[0] == 1 {
			num = true
		}
		v.numericPrereleases = append(v.numericPrereleases, num)
		data = data[1:]
	}
	v.builds = nil
	n, data = decodeInt(data)
	for i := 0; i < n; i++ {
		buff, data = decodeArray(data)
		v.builds = append(v.builds, buff)
	}
	return nil
}

// MarshalBinary implements json.Marshaler
func (v *RelaxedVersion) MarshalBinary() ([]byte, error) {
	res := new(bytes.Buffer)
	if len(v.customversion) > 0 {
		_, _ = res.Write([]byte{0})
		_, _ = res.Write(marshalByteArray(v.customversion))
		return res.Bytes(), nil
	}
	res.Write([]byte{1})
	d, _ := v.version.MarshalBinary() // can't fail
	_, _ = res.Write(d)
	return res.Bytes(), nil
}

// UnmarshalBinary implements json.Unmarshaler
func (v *RelaxedVersion) UnmarshalBinary(data []byte) error {
	if data[0] == 0 {
		v.customversion, _ = decodeArray(data[1:])
		v.version = nil
		return nil
	}

	v.customversion = nil
	v.version = &Version{}
	return v.version.UnmarshalBinary(data[1:])
}

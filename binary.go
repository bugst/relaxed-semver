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

// MarshalBinary implements binary custom encoding
func (v *Version) MarshalBinary() ([]byte, error) {
	// TODO could be preallocated without bytes.Buffer
	res := new(bytes.Buffer)
	intBuff := [4]byte{}
	_, _ = res.Write(marshalByteArray([]byte(v.raw)))
	binary.BigEndian.PutUint32(intBuff[:], uint32(v.major))
	_, _ = res.Write(intBuff[:])
	binary.BigEndian.PutUint32(intBuff[:], uint32(v.minor))
	_, _ = res.Write(intBuff[:])
	binary.BigEndian.PutUint32(intBuff[:], uint32(v.patch))
	_, _ = res.Write(intBuff[:])
	binary.BigEndian.PutUint32(intBuff[:], uint32(v.prerelease))
	_, _ = res.Write(intBuff[:])
	binary.BigEndian.PutUint32(intBuff[:], uint32(v.build))
	_, _ = res.Write(intBuff[:])
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

	buff, data = decodeArray(data)
	v.raw = string(buff)
	v.major, data = decodeInt(data)
	v.minor, data = decodeInt(data)
	v.patch, data = decodeInt(data)
	v.prerelease, data = decodeInt(data)
	v.build, _ = decodeInt(data)
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

// Copyright 2018-2022 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"testing"
)

var result int

func BenchmarkNumericArray(b *testing.B) {
	count := 0
	for n := 0; n < b.N; n++ {
		c := byte(n)
		if numeric[c] {
			count++
		}
	}
	result = count
}

func BenchmarkNumericFunction(b *testing.B) {
	count := 0
	for n := 0; n < b.N; n++ {
		c := byte(n)
		if isNumeric(c) {
			count++
		}
	}
	result = count
}

func BenchmarkIdentifierArray(b *testing.B) {
	count := 0
	for n := 0; n < b.N; n++ {
		c := byte(n)
		if identifier[c] {
			count++
		}
	}
	result = count
}

func BenchmarkIdentifierFunction(b *testing.B) {
	count := 0
	for n := 0; n < b.N; n++ {
		c := byte(n)
		if isIdentifier(c) {
			count++
		}
	}
	result = count
}

func BenchmarkVersionSeparatorArray(b *testing.B) {
	count := 0
	for n := 0; n < b.N; n++ {
		c := byte(n)
		if versionSeparator[c] {
			count++
		}
	}
	result = count
}

func BenchmarkVersionSeparatorFunction(b *testing.B) {
	count := 0
	for n := 0; n < b.N; n++ {
		c := byte(n)
		if isVersionSeparator(c) {
			count++
		}
	}
	result = count
}

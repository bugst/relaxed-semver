//
// Copyright 2018-2023 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"testing"
)

var list = []string{
	"0.0.1-rc.0",       // 0
	"0.0.1-rc.0+build", // 1
	"0.0.1-rc.1",       // 2
	"0.0.1",            // 3
	"0.0.1+build",      // 4
	"0.0.2-rc.1",       // 5 - BREAKING CHANGE
	"0.0.2-rc.1+build", // 6
	"0.0.2",            // 7
	"0.0.2+build",      // 8
	"0.0.3-rc.1",       // 9 - BREAKING CHANGE
	"0.0.3-rc.2",       // 10
	"0.0.3",            // 11
	"0.1.0",            // 12 - BREAKING CHANGE
	"0.3.3-rc.0",       // 13 - BREAKING CHANGE
	"0.3.3-rc.1",       // 14
	"0.3.3",            // 15
	"0.3.3+build",      // 16
	"0.3.4-rc.1",       // 17
	"0.3.4",            // 18
	"0.4.0",            // 19 - BREAKING CHANGE
	"1.0.0-rc",         // 20 - BREAKING CHANGE
	"1.0.0",            // 21
	"1.0.0+build",      // 22
	"1.2.1-rc",         // 23
	"1.2.1",            // 24
	"1.2.1+build",      // 25
	"1.2.3-rc.2",       // 26
	"1.2.3-rc.2+build", // 27
	"1.2.3",            // 28
	"1.2.3+build",      // 29
	"1.2.4",            // 30
	"1.3.0-rc.0+build", // 31
	"1.3.0",            // 32
	"1.3.0+build",      // 33
	"1.3.1-rc.0",       // 34
	"1.3.1-rc.1",       // 35
	"1.3.1",            // 36
	"1.3.5",            // 37
	"2.0.0-rc",         // 38 - BREAKING CHANGE
	"2.0.0-rc+build",   // 39
	"2.0.0",            // 40
	"2.0.0+build",      // 41
	"2.1.0-rc",         // 42
	"2.1.0-rc+build",   // 43
	"2.1.0",            // 44
	"2.1.0+build",      // 45
	"2.1.3-rc",         // 46
	"2.1.3",            // 47
	"2.3.0",            // 48
	"2.3.1",            // 49
	"3.0.0",            // 50 - BREAKING CHANGE
}

func BenchmarkVersionParser(b *testing.B) {
	res := &Version{}
	for i := 0; i < b.N; i++ {
		for _, v := range list {
			res.raw = v
			res.bytes = []byte(v)
			parse(res)
		}
	}

	// $ go test -benchmem -run=^$ -bench ^BenchmarkVersionParser$ go.bug.st/relaxed-semver
	// goos: linux
	// goarch: amd64
	// pkg: go.bug.st/relaxed-semver
	// cpu: AMD Ryzen 5 3600 6-Core Processor

	// Results for v0.11.0:
	// BenchmarkVersionParser-12    	  188611	      7715 ns/op	    8557 B/op	      51 allocs/op

	// Results for v0.12.0:  \o/
	// BenchmarkVersionParser-12    	  479626	      3719 ns/op	     616 B/op	      51 allocs/op
}

func BenchmarkVersionComparator(b *testing.B) {
	b.StopTimer()
	vList := []*Version{}
	for _, in := range list {
		vList = append(vList, MustParse(in))
	}
	l := len(vList)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		// cross compare all versions
		for x := 0; x < l; x++ {
			for y := 0; y < l; y++ {
				vList[x].CompareTo(vList[y])
			}
		}
	}

	// $ go test -benchmem -run=^$ -bench ^BenchmarkVersionComparator$ go.bug.st/relaxed-semver -v
	// goos: linux
	// goarch: amd64
	// pkg: go.bug.st/relaxed-semver
	// cpu: AMD Ryzen 5 3600 6-Core Processor

	// Results for v0.11.0:
	// BenchmarkVersionComparator-12    	   74793	     17347 ns/op	       0 B/op	       0 allocs/op

	// Results for v0.12.0:  :-D
	// BenchmarkVersionComparator-12    	  101772	     11720 ns/op	       0 B/op	       0 allocs/op
}

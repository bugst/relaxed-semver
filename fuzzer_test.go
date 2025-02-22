//
// Copyright 2018-2025 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"cmp"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func FuzzParser(f *testing.F) {
	f.Add("0-")
	f.Add("9+V.1.3.1.1.1.1.1.9.1.3.1.1.1.1.1.9.1")
	f.Add("757202302126572447535523755625377572023021265724475355237556000065.")
	f.Add("4+57023021265724475355237556253000065.\x00")
	f.Add("0--xAbbbbeAcAcCECaBaAAaAfeAdEBe-xfCCfBBfAEdcfFebeDxfCCfBBfAEdcfFebeD")
	f.Add("9-V.1.3.1.1.1.1.1.9.1")
	f.Add("3.1.12")
	f.Add("4.474368202171")
	f.Add("2-V.1.3.1.1.1.1.1.V.1.3.1.1.1.1.1.9.1")
	f.Add("4.02")
	f.Add("0-057e-33304e-91094BfAEd6cf6379282317958222700xfCCfB5BfAEd6cfFebe7D")
	f.Add("0\x01")
	f.Add("9+3.1.4.0.1.-.1.R.1")
	f.Add("1.1.1\x0e")
	f.Add("4-V.t.t.t.e.e.V.t.e.V.t.e.e.V.t.e.V")
	f.Add("1.0x")
	f.Add("1+1.1.1")
	f.Add("0-0-.0.0")
	f.Add("4.4.4740683-m")
	f.Add("0-0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0")
	f.Add("0-57e-33304e-910946379282317068958222700xfC792823CfB5BfAEd6cfFebe7D")
	f.Add("9298023223876953125\v")
	f.Add("4.4.474068")
	f.Add("1+1.1.1.3")
	f.Add("1.1.13682021700570577572023021265724475355237556253000065.")
	f.Add("0.0.")
	f.Add("0-9+0")
	f.Add("9-3.1.9.1.\x81")
	f.Add("3-9625")
	f.Add("4.0-0")
	f.Add("1+1")
	f.Add("1+1.1")
	f.Add("390625")
	f.Add("1.98023223876953125\v")
	f.Add("1+1.1.")
	f.Add("1.1.-")
	f.Add("1.1.0-0")
	f.Add("3-02")
	f.Add("1+1.1\x88")
	f.Add("0-0.0.0.0.0.0.0.0.00")
	f.Add("0\xab")
	f.Add("5.-")
	f.Add("0-0.0")
	f.Add("1.1.0")
	f.Add("9+3.1.1.1.-.1.R.1")
	f.Add("1.1.9+3")
	f.Add("1.1368202170057057757202302126572447535523755625377572023021265724475355237556000065.")
	f.Add("3-a-after ")
	f.Add("9-3.3.1.9.1.")
	f.Add("1+19446951953614188\x00")
	f.Add("1+")
	f.Add("1.1\x0e")
	f.Add("0-0x.0")
	f.Add("2-V.V.t.e.V")
	f.Add("0-0.0.0.0.0.00")
	f.Add("1+1.fterafter ")
	f.Add("1++")
	f.Add("0-9+")
	f.Add("1.1.11.")
	f.Add("12")
	f.Add("0-57e-33304e-910946379282317958222700xfC792823CfB5BfAEd6cfFebe7D910946379282317958222700xfC792823CfB5BfAEd6cfFebe7D")
	f.Add("3-a-fterafter ")
	f.Add("3-")
	f.Add("3-<")
	f.Add("3-4.9.4")
	f.Add("1+57e-3330-57e-3346379282317068958222700xfC792823CfB5BfAEd6cfFebe7D")
	f.Add("0-1690993073057962936658730400845563-0xacC.-0xFe34b9")
	f.Add("0-057e-33304e-910946379282317958222700xfCCfB5BfAEd6cfFebe7D")
	f.Add("4-V.t.t.V.t.t.t.e.e.V.t.e.V.t.e.e.V.t.e.t.e.e.V.t.e.V.t.e.e.V.t.e.V")
	f.Add("@")
	f.Add("9-3.1.9.1.n")
	f.Add("0-0.0-.0.0.0.0.0.0.0.0.0.0.0.0.0.0.00")
	f.Add("9+3.1.1.1.R.1")
	f.Add("9-3.1.1.1.9.1.")
	f.Add("9+3.1.3.1.1.1.1.1.V.1.3.1.1.1.1.1.3.1.1.1.1.1.V.1.3.1.1.1.1.1.9.1.9.1")
	f.Add("3.1.1-2")
	f.Add("2-V.t.e.V.t.e.V")
	f.Add("3-9402004353104906.474368202171")
	f.Add("3.1.8081828384858687888912")
	f.Add("1+1.1.1.1.")
	f.Add("9+3.1.4.0.1.-.1.R.1.")
	f.Add("9+3.1.1.1.R.1.")
	f.Add("0-0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0-.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0")
	f.Add("0")
	f.Add("0-057e-33304e-910946379282317958222700xfC792823CfB5BfAEd6cfFebe7D")
	f.Add("34694469519536141888238489627838134765\x00")
	f.Add("1.1.1")
	f.Add("0-0.0.0.0.0.0.00")
	f.Add("1+1.1.1.")
	f.Add("9+3.1.1.1.-.1.R.1.")
	f.Add("0.0.0.")
	f.Add("9-3.1.1.1.9.1")
	f.Add("9+3.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.-.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0")
	f.Add("3-after ")
	f.Add("2-V.t.t.e.e.V.t.e.V")
	f.Add("0-0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0")
	f.Add("1.1.1368202170057057757202302126572447535523755625377572023021265724475355237556000065.")
	f.Add("")
	f.Add("1.1.0\x0e")
	f.Add("4.4743682021710570577572023021265724475355237556253000065.02")
	f.Add("01")
	f.Add("9+V.1.3.1.1.1.1.1.9.1")
	f.Add("0.0.0.0")
	f.Add("1.0")
	f.Add("2-V.1.3.1.1.1.1.1.V.1.3.1.1.1.1.1.3.1.1.1.1.1.V.1.3.1.1.1.1.1.9.1.9.1")
	f.Add("3-m")
	f.Add("9\v")
	f.Add("2-V.V.e.V")
	f.Fuzz(func(t *testing.T, in string) {
		// ParseRelaxed should always succeed
		r := ParseRelaxed(in)
		if r.String() != in {
			t.Fatalf("reserialized relaxed string != deserialized string (in=%v)", in)
		}
		if r.CompareTo(r) != 0 {
			t.Fatalf("compare != 0 while comparing with self (in=%v)", in)
		}

		// Parse should succeed only if the input is a valid semver
		v, err := Parse(in)
		if err != nil {
			if v != nil {
				t.Fatalf("v != nil on error (in=%v)", in)
			}
			return
		}
		if v.String() != in {
			t.Fatalf("reserialized string != deserialized string (in=%v)", in)
		}
		if v.CompareTo(v) != 0 {
			t.Fatalf("compare != 0 while comparing with self (in=%v)", in)
		}
		v.Normalize()
		if v.CompareTo(v) != 0 {
			t.Fatalf("compare != 0 while comparing with self (in=%v)", in)
		}
	})
}

func FuzzComparators(f *testing.F) {
	f.Add("1.2.4", "0.0.1-rc.0")
	f.Add("1.3.0-rc.0+build", "0.0.1-rc.0+build")
	f.Add("1.3.0", "0.0.1-rc.1")
	f.Add("1.3.0+build", "0.0.1")
	f.Add("0.0.1+build", "0.0.1-rc.0")
	f.Add("0.0.2-rc.1", "0.0.1-rc.0+build")
	f.Add("0.0.2-rc.1+build", "0.0.1-rc.1")
	f.Add("0.0.2", "0.0.1")
	f.Add("0.0.2+build", "0.0.1+build")
	f.Add("0.0.3-rc.1", "0.0.2-rc.1")
	f.Add("0.0.3-rc.2", "0.0.2-rc.1+build")
	f.Add("0.0.3", "0.0.2")
	f.Add("0.1.0", "0.0.2+build")
	f.Add("0.3.3-rc.0", "0.0.3-rc.1")
	f.Add("0.3.3-rc.1", "0.0.3-rc.2")
	f.Add("0.3.3", "0.0.3")
	f.Add("0.3.3+build", "0.1.0")
	f.Add("0.3.4-rc.1", "0.3.3-rc.0")
	f.Add("0.3.4", "0.3.3-rc.1")
	f.Add("0.4.0", "0.3.3")
	f.Add("1.0.0-rc", "0.3.3+build")
	f.Add("1.0.0", "0.3.4-rc.1")
	f.Add("1.0.0+build", "0.3.4")
	f.Add("1.2.1-rc", "0.4.0")
	f.Add("1.2.1", "1.0.0-rc")
	f.Add("1.2.1+build", "1.0.0")
	f.Add("1.2.3-rc.2", "1.0.0+build")
	f.Add("1.2.3-rc.2+build", "1.2.1-rc")
	f.Add("1.2.3", "1.2.1")
	f.Add("1.2.3+build", "1.2.1+build")
	f.Add("1.2.4", "1.2.3-rc.2")
	f.Add("1.3.0-rc.0+build", "1.2.3-rc.2+build")
	f.Add("1.3.0", "1.2.3")
	f.Add("1.3.0+build", "1.2.3+build")
	f.Add("1.3.1-rc.0", "1.2.4")
	f.Add("1.3.1-rc.1", "1.3.0-rc.0+build")
	f.Add("1.3.1", "1.3.0")
	f.Add("1.3.5", "1.3.0+build")
	f.Add("2.0.0-rc", "1.3.1-rc.0")
	f.Add("2.0.0-rc+build", "1.3.1-rc.1")
	f.Add("2.0.0", "1.3.1")
	f.Add("2.0.0+build", "1.3.5")
	f.Add("2.1.0-rc", "2.0.0-rc")
	f.Add("2.1.0-rc+build", "2.0.0-rc+build")
	f.Add("2.1.0", "2.0.0")
	f.Add("2.1.0+build", "2.0.0+build")
	f.Add("2.1.3-rc", "2.1.0-rc")
	f.Add("2.1.3", "2.1.0-rc+build")
	f.Add("2.3.0", "2.1.0")
	f.Add("2.3.1", "2.1.0+build")
	f.Add("3.0.0", "2.1.3-rc")
	f.Fuzz(func(t *testing.T, a, b string) {
		va, err := Parse(a)
		if err != nil {
			return
		}
		vb, err := Parse(b)
		if err != nil {
			return
		}
		fmt.Println(va.SortableString(), vb.SortableString())
		require.Equal(t, va.CompareTo(vb), cmp.Compare(va.SortableString(), vb.SortableString()), "Comparing: %s and %s", a, b)
	})
}

//
// Copyright 2018-2025 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLDriverInterfaces(t *testing.T) {

	t.Run("Version", func(t *testing.T) {
		// Test Version Scan/Value
		v := &Version{}
		if _, ok := interface{}(v).(driver.Valuer); !ok {
			t.Error("Version does not implement driver.Valuer")
		}
		if _, ok := interface{}(v).(driver.Valuer); !ok {
			t.Error("Version does not implement driver.Valuer")
		}
		require.Error(t, v.Scan(1))
		require.Error(t, v.Scan(nil))
		require.Error(t, v.Scan("123asdf"))
		require.NoError(t, v.Scan("1.2.3-rc.1+build.2"))
		require.Equal(t, "1.2.3-rc.1+build.2", v.String())
		d, err := v.Value()
		require.NoError(t, err)
		require.Equal(t, "1.2.3-rc.1+build.2", d)
	})

	t.Run("RelaxedVersion", func(t *testing.T) {
		// Test RelaxedVersion Scan/Value
		rv := &RelaxedVersion{}
		if _, ok := interface{}(rv).(driver.Valuer); !ok {
			t.Error("RelaxedVersion does not implement driver.Valuer")
		}
		if _, ok := interface{}(rv).(driver.Valuer); !ok {
			t.Error("RelaxedVersion does not implement driver.Valuer")
		}
		require.Error(t, rv.Scan(1))
		require.Error(t, rv.Scan(nil))
		require.NoError(t, rv.Scan("4.5.6-rc.1+build.2"))
		require.Empty(t, rv.customversion)
		require.NotNil(t, rv.version)
		require.Equal(t, "4.5.6-rc.1+build.2", rv.String())
		rd, err := rv.Value()
		require.NoError(t, err)
		require.Equal(t, "4.5.6-rc.1+build.2", rd)

		require.NoError(t, rv.Scan("a1-2.2-3.3"))
		require.NotEmpty(t, rv.customversion)
		require.Nil(t, rv.version)
		require.Equal(t, "a1-2.2-3.3", rv.String())
		rd2, err := rv.Value()
		require.NoError(t, err)
		require.Equal(t, "a1-2.2-3.3", rd2)
	})
}

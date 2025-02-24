//
// Copyright 2018-2025 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"database/sql/driver"
	"fmt"
)

// Scan implements the sql.Scanner interface
func (v *Version) Scan(value interface{}) error {
	raw, ok := value.(string)
	if !ok {
		return fmt.Errorf("incompatible type %T for Version", value)
	}

	v.raw = raw
	v.bytes = []byte(v.raw)
	if err := parse(v); err != nil {
		return err
	}
	return nil
}

// Value implements the driver.Valuer interface
func (v *Version) Value() (driver.Value, error) {
	return v.raw, nil
}

// Scan implements the sql.Scanner interface
func (v *RelaxedVersion) Scan(value interface{}) error {
	raw, ok := value.(string)
	if !ok {
		return fmt.Errorf("incompatible type %T for Version", value)
	}

	res := ParseRelaxed(raw)
	*v = *res
	return nil
}

// Value implements the driver.Valuer interface
func (v *RelaxedVersion) Value() (driver.Value, error) {
	if v.version != nil {
		return v.version.raw, nil
	}
	return string(v.customversion), nil
}

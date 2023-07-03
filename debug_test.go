//
// Copyright 2018-2023 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

import (
	"fmt"
	rtdebug "runtime/debug"
	"strings"
	"testing"
)

func init() {
	debug = func(format string, a ...interface{}) {
		level := strings.Count(string(rtdebug.Stack()), "\n")
		for i := 0; i < level; i++ {
			fmt.Print(" ")
		}
		if a != nil {
			fmt.Printf(format, a...)
			fmt.Println()
		} else {
			fmt.Println(format)
		}
	}
}

func TestNoopDebug(t *testing.T) {
	noopDebug("just for coverage!")
}

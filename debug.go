//
// Copyright 2018-2025 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package semver

var debug = noopDebug

func noopDebug(format string, a ...interface{}) {}

// Copyright (c) 2020 by meng.  All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package proxy

import "testing"

// TestCheckHeath
func TestCheckHeath(t *testing.T) {
	var backend map[string]*BackendEnd

	backend = make(map[string]*BackendEnd)
	backend["127.0.0.1:83"] = &BackendEnd{SvrStr: "127.0.0.1:83", IsUp: false, FailTimes: 1, RiseTimes: 1}
	backend["127.0.0.1:84"] = &BackendEnd{SvrStr: "127.0.0.1:84", IsUp: true, FailTimes: 1, RiseTimes: 1}

	checkHeath(backend)
}

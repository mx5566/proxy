// Copyright (c) 2020 by meng.  All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package proxy

import "testing"

func TestParseConfigFile(t *testing.T) {

	path := "./logs/tcpProxy.log"
	_ = parseConfigFile(path)
}

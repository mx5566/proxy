// Copyright (c) 2020 by meng.  All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package proxy

import "testing"

// go test -v log.go log_test.go config.go
func TestInitLog(t *testing.T) {

	var config LogConfig
	config.Level = 0
	config.ServiceName = "TestInitLog"
	config.Compress = true
	config.MaxAge = 7
	config.MaxSize = 10
	config.MaxBackup = 10

	initLog(&config)
}

// Copyright (c) 2020 by meng.  All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"testing"
)

func TestProxy_InitProxy(t *testing.T) {
	var config ProxyConfig
	config.Bind = "0.0.0.0:9999"
	config.Stats = "0.0.0.0:19090"

	var logConfig LogConfig

	logConfig.Level = 0
	logConfig.ServiceName = "TestInitLog"
	logConfig.Compress = true
	logConfig.MaxAge = 7
	logConfig.MaxSize = 10
	logConfig.MaxBackup = 10

	proxy.InitProxy(&config)
}

// go test -v proxy.go proxy_test.go config.go log.go

func TestCreateProxy(t *testing.T) {
	CreateProxy("./config.yaml")
}

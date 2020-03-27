package proxy

import (
	"testing"
)

func TestProxy_InitProxy(t *testing.T) {
	var config ProxyConfig
	config.Bind = "0.0.0.0:9999"
	config.WaitQueueLen = 100
	config.MaxConn = 50
	config.Timeout = 5
	config.FailOver = 3
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

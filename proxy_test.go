package proxy

import (
	"fmt"
	"os"
	"testing"
)

var proxy Proxy

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
func TestMain(m *testing.M) {
	err := parseConfigFile("./config.yaml")
	if err != nil {
		panic("load config file error" + err.Error())
		return
	}

	// 日志模块
	initLog(&config.Log)

	fmt.Println(config)

	var proxy Proxy
	// 初始化后端服务器
	proxy.InitBackEnd(&config)
	// 初始化代理模块
	proxy.InitProxy(&config)

	// 需要状态检查服务器

	os.Exit(m.Run())

}

func TestCreateProxy(t *testing.T) {
	CreateProxy("./config.yaml")
}

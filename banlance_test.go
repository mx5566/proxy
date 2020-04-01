package proxy

import (
	"fmt"
	"net"
	"testing"
)

func TestBanlance(t *testing.T) {
	var banlance Banlance

	backend := []string{"127.0.0.1:8080"}

	heatch := HeatchConfig{
		Interval:             3000,
		Rise:                 3,
		Fall:                 3,
		Timeout:              1000,
		Type:                 "http",
		DefaultDown:          false,
		CheckHttpSend:        "HEAD / HTTP/1.1\r\n\r\n",
		CheckHttpExceptAlive: []string{"http_2xx", "http_3xx"},
		Port:                 8080,
	}

	// 初始化banlance
	banlance.Init(backend, heatch)

	conn, err := net.Dial("tcp", backend[0])

	if conn != nil {
		fmt.Println(banlance.GetBackEndServer(conn))

	} else {
		fmt.Println(err)
	}

	// 服务治理
	stats := []BackStats{{
		IsUp:   true,
		SvrStr: "127.0.0.1:8080",
	}}
	banlance.ServiceGovernance(stats)

}

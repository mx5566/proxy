// Copyright (c) 2020 by meng.  All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"net"

	"stathat.com/c/consistent"
)

type Banlance struct {
	// 后端服务器的配置
	backend map[string]*BackendEnd

	// 对后端IP进行一致性哈希到环对应的位置
	pConsisthash *consistent.Consistent
}

func (this *Banlance) Init(backend []string, heatch HeatchConfig) {
	logger.Info("init banlance")
	for _, svr := range backend {
		// 把对应的后端服务器加入到哈希环
		this.pConsisthash.Add(svr)

		logger.Info(svr)
		this.backend[svr] = &BackendEnd{
			SvrStr:    svr,
			IsUp:      !heatch.DefaultDown,
			FailTimes: 0,
			RiseTimes: 0,
		}
	}
}

// 获得后端服务器
func (this *Banlance) GetBackEndServer(conn net.Conn) string {
	// 127.0.0.1:80
	clientIp := conn.RemoteAddr().String()

	// 从哈希环获取对应的哈希值也就是存入的Ip地址
	server, err := this.pConsisthash.Get(clientIp)

	if err != nil {
		return ""
	}

	svr, ok := this.backend[server]

	if !ok {
		return ""
	}

	return svr.SvrStr
}

// 后端应用服务器服务治理
func (this *Banlance) ServiceGovernance(stats []BackStats) {

	for _, stat := range stats {
		// 负载服务器更新
		if stat.IsUp == true {
			this.pConsisthash.Add(stat.SvrStr)
		} else {
			this.pConsisthash.Remove(stat.SvrStr)
		}
	}
}

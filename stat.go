// Copyright (c) 2020 by meng.  All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"net/http"
	"sync"
)

// Stat
// Real Server States
type Stat struct {
	// 客户端连接池--
	pool sync.Pool

	// http
	mux *http.ServeMux
}

func (this *Stat) StartStat() {
	logger.Info("start stats keep")

	// 初始化池子的对象构造
	this.pool = sync.Pool{
		New: func() interface{} {
			return &Context{Request: nil, ResponseWriter: nil}
		},
	}
	this.mux = http.NewServeMux()

	// 注册处理函数
	this.RegisterRoute("/stats", StatsHandler)

	// 监听状态端口
	_ = http.ListenAndServe(config.Stats, this.mux)

}

// 注册路由处理函数
func (this *Stat) RegisterRoute(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	this.mux.HandleFunc(uri, f)
}

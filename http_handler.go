// Copyright (c) 2020 by meng.  All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package proxy

import "net/http"

// 监控状态处理函数
// 所有后端服务器的状态
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	context := proxy.stat.pool.Get().(*Context)
	context.Reset(w, r)

	defer proxy.stat.pool.Put(context)

	// 返回数据给客户端
	ret := retJson{}
	ret.Code = 0
	ret.Msg = "ok"

	var data map[string]interface{}
	data = make(map[string]interface{})
	for _, v := range proxy.banlance.backend {
		data[v.SvrStr] = *v
	}

	ret.Data = data

	context.ServerJson(data)
}

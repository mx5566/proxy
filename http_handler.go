package proxy

import "net/http"

// 监控状态处理函数
// 所有后端服务器的状态
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	context := proxy.pool.Get().(*Context)
	context.Reset(w, r)

	defer proxy.pool.Put(context)

	// 返回数据给客户端
	ret := retJson{}
	ret.Code = 0
	ret.Msg = "ok"

	var data map[string]interface{}
	data = make(map[string]interface{})
	for _, v := range proxy.backend {
		data[v.SvrStr] = *v
	}

	ret.Data = data

	context.ServerJson(data)
}

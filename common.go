package proxy

// 返回的状态数据格式
type retJson struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
}

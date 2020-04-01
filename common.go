package proxy

// 返回的状态数据格式
type retJson struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
}

// BackendSvr Type
type BackendEnd struct {
	SvrStr    string `json:"svrStr"`
	IsUp      bool   `json:"isUp"`      // is Up or Down
	FailTimes int    `json:"failTimes"` // 失败次数
	RiseTimes int    `json:"riseTimes"` // 连接成功的次数
}

type BackStats struct {
	IsUp   bool // is Up or Down
	SvrStr string
}

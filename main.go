package proxy

var config ProxyConfig

func main() {

	err := parseConfigFile("./config.yaml")
	if err != nil {
		panic("load config file error" + err.Error())
		return
	}

	// 日志模块
	initLog(&config.Log)

	var proxy Proxy
	// 初始化后端服务器
	proxy.InitBackEnd(&config)
	// 初始化代理模块
	proxy.InitProxy(&config)

	// 需要状态检查服务器
}

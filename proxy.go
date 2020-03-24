package main

import (
	"net"
	"time"

	"stathat.com/c/consistent"
)

type Proxy struct {
	// 后端服务器的配置
	backend map[string]*BackendEnd

	// 对后端IP进行一致性哈希到环对应的位置
	pConsisthash *consistent.Consistent
}

// BackendSvr Type
type BackendEnd struct {
	svrStr    string
	isUp      bool // is Up or Down
	failTimes int
}

func (this *Proxy) InitProxy(proxyConfig *ProxyConfig) {
	// 启动代理服务器
	listener, err := net.Listen("tcp", proxyConfig.Bind)
	if err != nil {
		//panic("init proxy error")
		logger.Error(err.Error())
		return
	}

	logger.Info("init proxy listen ok")

	for {
		// 代理服务器监听客户端的连接
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("accept errro " + err.Error())
			break
		}

		// 处理客户端连接
		go func(conn net.Conn) {
			this.HandleConnect(conn)
			// 连接处理完了
			logger.Info("conn handle ok")
		}(conn)
	}

}

func (this *Proxy) HandleConnect(conn net.Conn) {
	defer conn.Close()

	// 把客户端的连接转发给后端服务器
	// 获取一个后端服务器
	server := this.GetBackEnd(conn)
	if server == "" {
		return
	}

	// 找到后端服务器 通过代理连接后端服务器
	// 连接成功转发客户端的数据到后端服务器连接
	serverSession, err := net.Dial("tcp", server)
	if err != nil {
		logger.Info("connect back server error " + server + " " + err.Error())
		return
	}

	var ok chan bool

	ok = make(chan bool, 2)
	// 把客户端的数据-写到服务器端
	go this.Copy(conn, serverSession, ok)
	// 把服务器端的数据写到客户端
	go this.Copy(serverSession, conn, ok)

	// 通过管道控制
	<-ok
	<-ok

	// 说明两端出错了
	serverSession.Close()
}

func (this *Proxy) Copy(from net.Conn, to net.Conn, ok chan bool) {
	var err error
	var read int
	bytes := make([]byte, 256)

	for {
		select {
		default:
			_ = from.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
			read, err = from.Read(bytes)
			if err != nil {
				ok <- true
				return
			}

			_ = to.SetWriteDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
			_, err = to.Write(bytes[:read])
			if err != nil {
				ok <- true
				return
			}
		}
	}
}

func (this *Proxy) InitBackEnd(proxyConfig *ProxyConfig) {
	this.pConsisthash = consistent.New()

	for _, svr := range proxyConfig.Backend {
		// 把对应的后端服务器加入到哈希环
		this.pConsisthash.Add(svr)

		this.backend[svr] = &BackendEnd{
			svrStr:    svr,
			isUp:      true,
			failTimes: 0,
		}
	}
}

func (this *Proxy) GetBackEnd(conn net.Conn) string {
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

	return svr.svrStr
}

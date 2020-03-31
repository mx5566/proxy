package proxy

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"fmt"

	"stathat.com/c/consistent"
)

// 初始化全局proxy
var proxy Proxy

type Proxy struct {
	// 后端服务器的配置
	backend map[string]*BackendEnd

	// 对后端IP进行一致性哈希到环对应的位置
	pConsisthash *consistent.Consistent

	// 客户端连接池--
	pool sync.Pool

	mux *http.ServeMux

	isShutdown chan bool

	limiter LimitInterface
}

func (this *Proxy) InitProxy(proxyConfig *ProxyConfig) {

	// 启动代理服务器
	listener, err := net.Listen("tcp", proxyConfig.Bind)
	if err != nil {
		//panic("init proxy error")
		logger.Error(err.Error())
		return
	}

	defer listener.Close()

	logger.Info("init proxy listen ok")

	// 接受关闭数据 1
	this.onSignalClose(listener)

	// 创建限流器 2
	this.limiter = NewQueueLimter(proxyConfig.Limter.WaitQueueLen, proxyConfig.Limter.MaxConn)
	if this.limiter == nil {
		return
	}

	// 启动限流器处理流程
	this.limiter.Bind(this.HandleConnect)

close:
	for {
		// 代理服务器监听客户端的连接
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("accept errro " + err.Error())
			break close
		}

		// 限流器连接添加 3
		if this.limiter.IsAvalivale() {
			this.limiter.SetWaitQueue(conn)
		} else {
			_, _ = conn.Write([]byte("server full"))
			conn.Close()
		}
	}

}

func (this *Proxy) onSignalClose(listener net.Listener) {
	go func(ln net.Listener) {
		<-this.isShutdown
		if err := ln.Close(); err != nil {
			logger.Error(err.Error())
		}
	}(listener)
}

func (this *Proxy) HandleConnect(conn interface{}) {
	connTemp := conn.(net.Conn)

	defer connTemp.Close()

	// 把客户端的连接转发给后端服务器
	// 获取一个后端服务器
	server := this.GetBackEnd(connTemp)
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
	go this.Copy(connTemp, serverSession, ok)
	// 把服务器端的数据写到客户端
	go this.Copy(serverSession, connTemp, ok)

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
			_ = from.SetReadDeadline(time.Now().Add(time.Duration(config.Heatch.Timeout) * time.Millisecond))
			read, err = from.Read(bytes)
			if err != nil {
				ok <- true
				return
			}

			_ = to.SetWriteDeadline(time.Now().Add(time.Duration(config.Heatch.Timeout) * time.Millisecond))
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
	this.backend = make(map[string]*BackendEnd)

	for _, svr := range proxyConfig.Backend {
		// 把对应的后端服务器加入到哈希环
		this.pConsisthash.Add(svr)

		logger.Info(svr)
		this.backend[svr] = &BackendEnd{
			SvrStr:    svr,
			IsUp:      !proxyConfig.Heatch.DefaultDown,
			FailTimes: 0,
			RiseTimes: 0,
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

	return svr.SvrStr
}

func (this *Proxy) CheckBackEnd() {
	go checkHeath(this.backend)
}

func (this *Proxy) StatsBackEnd() {
	logger.Info("start stats keep")

	go func() {
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

	}()
}

func (this *Proxy) RegisterRoute(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	this.mux.HandleFunc(uri, f)
}

// 信号处理
func (this *Proxy) OnSignalExit() {
	// 初始化管道
	this.isShutdown = make(chan bool, 1)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

		logger.Info("Wait OnSignalExit")
		sig := <-c

		pid := syscall.Getpid()

		// 通知主线程关闭服务器
		switch sig {
		case syscall.SIGHUP:
			logger.Info("syscall.SIGHUP")
			this.isShutdown <- true

		case syscall.SIGINT:
			str := fmt.Sprintf("%d ", pid)
			logger.Info(str + "Received SIGINT.")
			this.isShutdown <- true

		case syscall.SIGTERM:
			logger.Info(string(pid) + "Received SIGTERM.")
			this.isShutdown <- true
		default:
			str := fmt.Sprintf("Received %s: nothing i care about", sig)
			logger.Info(str)
		}

	}()
}

func CreateProxy(configPath string) {
	//err := parseConfigFile("./config.yaml")
	err := parseConfigFile(configPath)
	if err != nil {
		panic("load config file error" + err.Error())
		return
	}

	// 日志模块
	initLog(&config.Log)

	// 初始化后端服务器
	proxy.InitBackEnd(&config)

	// 后端检测
	proxy.CheckBackEnd()

	// 后端服务器的状态用来显示
	proxy.StatsBackEnd()

	proxy.OnSignalExit()

	// 初始化代理模块
	proxy.InitProxy(&config)

}

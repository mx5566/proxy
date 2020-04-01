package proxy

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"

	"stathat.com/c/consistent"
)

// 初始化全局proxy
var proxy Proxy

type Proxy struct {

	// 服务器关闭管道
	isShutdown chan bool

	// 限流器
	limiter LimitInterface

	// 负载均衡器
	banlance Banlance

	// 服务器查看接口
	stat Stat
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
	server := this.banlance.GetBackEndServer(connTemp)
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

func (this *Proxy) StartBanlance(proxyConfig *ProxyConfig) {
	this.banlance = Banlance{
		backend:      make(map[string]*BackendEnd),
		pConsisthash: consistent.New(),
	}

	// 初始化负载均衡对象
	this.banlance.Init(proxyConfig.Backend, proxyConfig.Heatch)
	// 健康检查
	go checkHeath(this.banlance.backend)
}

func (this *Proxy) StartStat() {
	// 创建状态类
	this.stat = Stat{mux: http.NewServeMux()}

	// 启动启动函数
	go this.stat.StartStat()
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
	proxy.StartBanlance(&config)

	// 后端服务器的状态用来显示
	proxy.StartStat()

	proxy.OnSignalExit()

	// 初始化代理模块
	proxy.InitProxy(&config)

}

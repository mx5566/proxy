package proxy

import (
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// HeatchMontior
type HeathMontior struct {
}

// TcpCheck
func (this *HeathMontior) TcpCheck(hConfig HeatchConfig, backend map[string]*BackendEnd) {

	// 定时循环检测 相关的参数  健康检查
	second := time.Duration(hConfig.Interval) * time.Millisecond
	t := time.Tick(second)
	for _ = range t {
		for _, v := range backend {
			conn, err := net.DialTimeout("tcp", v.svrStr, time.Duration(hConfig.Timeout/1000))
			if err != nil {
				v.failTimes++
				v.riseTimes = 0
				if v.failTimes > hConfig.Fall {
					v.isUp = false
				}
				continue
			}

			v.riseTimes++
			v.failTimes = 0
			if v.riseTimes >= hConfig.Rise {
				v.isUp = true
			}

			defer conn.Close()
		}
	}
}

// HttpCheck
func (this *HeathMontior) HttpCheck(hConfig HeatchConfig, backend map[string]*BackendEnd) {

	// 定时循环检测 相关的参数  健康检查
	second := time.Duration(hConfig.Interval) * time.Millisecond
	t := time.Tick(second)
	for _ = range t {
		for _, v := range backend {
			resp, err := http.Head(v.svrStr)
			//resp, err := http.Get(v.svrStr)
			if err != nil {
				v.failTimes++
				v.riseTimes = 0
				if v.failTimes > hConfig.Fall {
					v.isUp = false
				}
				continue
			}

			v.riseTimes++
			v.failTimes = 0
			if v.riseTimes >= hConfig.Rise {
				v.isUp = true
			}
		}
	}
}

// PortCheck
func (this *HeathMontior) PortCheck(hConfig HeatchConfig, backend map[string]*BackendEnd) {

}

// checkHeath heath check
func checkHeath(backend map[string]*BackendEnd) {
	heath := config.Heatch

	var montior HeathMontior

	// 类型检测
	switch heath.Type {
	case "tcp":
		logger.Info("heatch tcp")
		montior.TcpCheck(heath, backend)
	case "http":
		logger.Info("heatch http")
		montior.HttpCheck(heath, backend)
	case "port":
		logger.Info("heatch port")
		montior.PortCheck(heath, backend)
	default:
		logger.Error("not support heatch type", zap.String("heatchType", heath.Type))
	}
}

package proxy

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// HeatchMontior
type HeathMontior struct {
}

func (this *HeathMontior) ParseIP(port int, svrStr string) string {
	ipPort := strings.Split(svrStr, ":")

	length := len(ipPort)
	var requestStr string
	if length <= 0 || port == 0 {
		requestStr = svrStr
	} else {
		requestStr = fmt.Sprintf("%s:%d", ipPort[0], port)
	}

	return requestStr
}

// TcpCheck
func (this *HeathMontior) TcpCheck(hConfig HeatchConfig, backend map[string]*BackendEnd) {

	// 定时循环检测 相关的参数  健康检查
	second := time.Duration(hConfig.Interval) * time.Millisecond
	t := time.Tick(second)
	for _ = range t {
		for _, v := range backend {
			requestStr := this.ParseIP(hConfig.Port, v.SvrStr)

			timeout := time.Duration(hConfig.Timeout) * time.Millisecond
			conn, err := net.DialTimeout("tcp", requestStr, timeout)
			if err != nil {
				v.FailTimes++
				v.RiseTimes = 0
				if v.FailTimes > hConfig.Fall {
					v.IsUp = false
				}
				continue
			}

			v.RiseTimes++
			v.FailTimes = 0
			if v.RiseTimes >= hConfig.Rise {
				v.IsUp = true
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
			requestStr := this.ParseIP(hConfig.Port, v.SvrStr)

			resp, err := http.Head(requestStr)
			if err != nil {
				v.FailTimes++
				v.RiseTimes = 0
				if v.FailTimes > hConfig.Fall {
					v.IsUp = false
				}
				continue
			}

			// 判断状态码是不是check_http_expect_alive
			// http_2xx
			// http_3xx
			var find = false
			for _, value := range hConfig.CheckHttpExceptAlive {
				alive := "http_" + string(strconv.Itoa(resp.StatusCode)[0]) + "xx"
				if alive == value {
					find = true
					break
				}
			}

			if !find {
				v.FailTimes++
				v.RiseTimes = 0
				if v.FailTimes > hConfig.Fall {
					v.IsUp = false
				}

				continue
			}

			v.RiseTimes++
			v.FailTimes = 0
			if v.RiseTimes >= hConfig.Rise {
				v.IsUp = true
			}
		}
	}
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
	default:
		logger.Error("not support heatch type", zap.String("heatchType", heath.Type))
	}
}

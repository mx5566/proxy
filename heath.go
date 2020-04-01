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

		stats := []BackStats{}
		for _, v := range backend {

			// 闭包函数
			fill := func() {
				stats = append(stats, BackStats{
					IsUp:   v.IsUp,
					SvrStr: v.SvrStr,
				})
			}

			beforeStats := v.IsUp
			requestStr := this.ParseIP(hConfig.Port, v.SvrStr)

			timeout := time.Duration(hConfig.Timeout) * time.Millisecond
			conn, err := net.DialTimeout("tcp", requestStr, timeout)
			if err != nil {
				v.FailTimes++
				v.RiseTimes = 0
				if v.FailTimes > hConfig.Fall && v.IsUp == true {
					v.IsUp = false
					// 服务器状态变了，需要更新负载均衡服务器的一致性
					fill()
				}

				continue
			}

			v.RiseTimes++
			v.FailTimes = 0
			if v.RiseTimes >= hConfig.Rise && v.IsUp == false {
				v.IsUp = true
				fill()
			}

			defer conn.Close()
		}

		// 通知到负载均衡器
		proxy.banlance.ServiceGovernance(stats)

	}
}

// HttpCheck
func (this *HeathMontior) HttpCheck(hConfig HeatchConfig, backend map[string]*BackendEnd) {

	// 定时循环检测 相关的参数  健康检查
	second := time.Duration(hConfig.Interval) * time.Millisecond
	t := time.Tick(second)
	for _ = range t {
		stats := []BackStats{}

		for _, v := range backend {

			// 闭包函数
			fill := func() {
				stats = append(stats, BackStats{
					IsUp:   v.IsUp,
					SvrStr: v.SvrStr,
				})
			}

			beforeStats := v.IsUp
			requestStr := this.ParseIP(hConfig.Port, v.SvrStr)
			resp, err := http.Head("http://" + requestStr)
			if err != nil {
				v.FailTimes++
				v.RiseTimes = 0
				if v.FailTimes > hConfig.Fall && v.IsUp == true {
					v.IsUp = false
					fill()
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
				if v.FailTimes > hConfig.Fall && v.IsUp == true {
					v.IsUp = false
					fill()
				}
				continue
			}

			v.RiseTimes++
			v.FailTimes = 0
			if v.RiseTimes >= hConfig.Rise && v.IsUp == false {
				v.IsUp = true
				fill()
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

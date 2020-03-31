# proxy
 tcpproxy

#可以通过下面的命令下载
```batch
go get github.com/mx5566/proxy

```

1. tcp代理
2. 采用一致性哈希算法进行负载均衡
3. 后端服务器心跳检查有tcp检查或者http检查、支持失败次数、成功次数、连接超时、允许设置初始的服务器状态、修改服务器心跳检测的端口 如果设置为0就采用后端服务器端口
4. 可以通过stats 获取后端服务器的状态列表
5. 日志采用的是uber的zap
6. 增加限流器 总共有三个


限流器
```batch
const (
	queueLimter        LimitType = iota + 1 // 队列模式限流
	tokenBucketLimter                       // golang官方库  golang.org/x/time/rate  bucket
	slideWindowLimiter                      // 滑动窗口限流器 tcp滑动窗口 -- 没有实现
	leakBucketLimter                        // 漏斗桶限流器
)



1对于队列模式queueLimter
采用的是二级模式
配置文件中
  wait_queue_len: 100 # 代表可以处于等待处理的队列
  max_conn: 50000 # 代表同时处于处理的连接 并发的用户

函数分为三步
-- NewQueueLimter(wait_queue_len, max_conn)
-- Bind(handler func(conn interface{}))
-- 		if this.limiter.IsAvalivale() {
   			this.limiter.SetWaitQueue(conn)
   		} 



2 tokenBucketLimter
  duration: 8  # 单位毫秒--速率
  captity: 100 # 容量
函数分为三步
-- NewTokenBucketLimiter(duration, captity)
-- Bind(handler func(conn interface{}))
-- 		if this.limiter.IsAvalivale() {
   			this.limiter.SetWaitQueue(conn)
   		} 

3 LeakBucketLimiter
  duration: 8  # 单位毫秒--速率
  captity: 100 # 容量
  name: "Test" # 限流器名字

函数分为三步
-- NewLeakBucketLimiter(name, captity, duration)
-- Bind(handler func(conn interface{}))
-- 		if this.limiter.IsAvalivale() {
   			this.limiter.SetWaitQueue(conn)
   		} 



```

例子直接通过调用下面的函数就可以
里面的yaml文件目录可以自己设定
```
    CreateProxy("./config.yaml")
```


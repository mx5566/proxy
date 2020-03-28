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

例子直接通过调用下面的函数就可以
里面的yaml文件目录可以自己设定
```
    CreateProxy("./config.yaml")
```


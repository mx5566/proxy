// Copyright (c) 2020 by meng.  All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"fmt"
	"time"

	"github.com/mx5566/proxy/leakybucket"

	"context"

	"golang.org/x/time/rate"
)

type LimitType int

const (
	queueLimter        LimitType = iota + 1 // 队列模式限流
	tokenBucketLimter                       // golang官方库  golang.org/x/time/rate  bucket
	slideWindowLimiter                      // 滑动窗口限流器 tcp滑动窗口
	leakBucketLimter                        // 漏斗桶限流器
)

// 限流接口
type LimitInterface interface {
	IsAvalivale() bool
	Bind(handler func(conn interface{}))
	SetWaitQueue(conn interface{})
}

///////////////////////////////QueueLimiter/////////////////////////////
// 通过队列实现限流
type QueueLimiter struct {
	waitQueue        chan interface{} // 等待队列类似c网络 listen的backlog
	availPools       chan bool        // 并发连接数
	initWaitQueueLen int
	initAvailConn    int
}

// NewQueueLimter
/**
	waitLength-最大的等待处理的长度
	maxConn-最大的并发处理长度

**/
func NewQueueLimter(waitLength, maxConn int) *QueueLimiter {
	limiter := &QueueLimiter{
		waitQueue:        make(chan interface{}, waitLength),
		availPools:       make(chan bool, maxConn),
		initWaitQueueLen: waitLength,
		initAvailConn:    maxConn,
	}

	// 预先初始化队列的长度，可用队列长度
	for i := 0; i < maxConn; i++ {
		limiter.availPools <- true
	}

	return limiter
}

// 等待队列是否还有空位
func (this QueueLimiter) IsAvalivale() bool {
	length := len(this.waitQueue)

	// 超过了等待队列的数量 说明超过了最大并发持续进行中
	if length >= this.initWaitQueueLen {
		return false
	}

	return true
}

// 限流器增加计数
func (this *QueueLimiter) SetWaitQueue(conn interface{}) {
	this.waitQueue <- conn
}

func (this QueueLimiter) Bind(handler func(conn interface{})) {
	go func() {
		for connection := range this.waitQueue {
			<-this.availPools
			go func(connection interface{}) {
				handler(connection)
				this.availPools <- true
				//logger.Info("conn handle ok on QueueLimiter")
			}(connection)
		}
	}()
}

///////////////////////////////QueueLimiter/////////////////////////////

//////////////////////////////TokenBucketLimiter/////////////////////////////

// 令牌桶实现限流
type TokenBucketLimiter struct {
	limiter *rate.Limiter // 令牌桶算法类
	handler func(conn interface{})
}

func NewTokenBucketLimiter(t int, token int) *TokenBucketLimiter {
	// 800ms
	limit := rate.Every(time.Duration(t) * time.Millisecond)

	bucket := &TokenBucketLimiter{limiter: rate.NewLimiter(limit, token)}
	return bucket
}

func (this *TokenBucketLimiter) IsAvalivale() bool {
	ctx, cancel := context.WithCancel(context.Background())

	// 取消cancel
	defer cancel()

	if err := this.limiter.Wait(ctx); err != nil {
		logger.Panic(err.Error())
	}

	return true
}

//
// bind handler function to  handler
func (this *TokenBucketLimiter) Bind(handler func(conn interface{})) {
	// TODO:
	// None
	this.handler = handler
}

// call handler function by conn
func (this *TokenBucketLimiter) SetWaitQueue(conn interface{}) {
	// TODO:
	// None
	go func(connection interface{}) {
		this.handler(connection)
		logger.Info("conn handle ok on BucketLimiter")

	}(conn)
}

//////////////////////////////TokenBucketLimiter/////////////////////////////

//////////////////////////////LeakBucketLimiter/////////////////////////////
// 漏斗桶限流算法
type LeakBucketLimiter struct {
	limiter leakybucket.BucketI
	handler func(conn interface{})
}

func NewLeakBucketLimiter(name string, cap uint, t int) *LeakBucketLimiter {
	bucketI, _ := leakybucket.New().Create(name, cap, time.Duration(t)*time.Millisecond)

	limit := &LeakBucketLimiter{limiter: bucketI}

	return limit
}

func (this *LeakBucketLimiter) IsAvalivale() bool {
	if _, err := this.limiter.Add(1); err != nil {
		return false
	}

	return true
}

//
// bind handler function to  handler
func (this *LeakBucketLimiter) Bind(handler func(conn interface{})) {
	// TODO:
	// None
	this.handler = handler
}

// call handler function by conn
func (this *LeakBucketLimiter) SetWaitQueue(conn interface{}) {
	// TODO:
	// None
	go func(connection interface{}) {
		this.handler(connection)
		logger.Info("conn handle ok on BucketLimiter")
	}(conn)
}

//////////////////////////////LeakBucketLimiter/////////////////////////////

// 不同的限流器的初始化接口
func CreateLimiter(lConfig LimiterConfig) LimitInterface {
	t := LimitType(lConfig.Type)
	var limiter LimitInterface
	switch t {
	case queueLimter:
		logger.Info("queueLimter")
		limiter = NewQueueLimter(lConfig.WaitQueueLen, lConfig.MaxConn)
	case tokenBucketLimter:
		logger.Info("tokenBucketLimter")
		limiter = NewTokenBucketLimiter(lConfig.Duration, int(lConfig.Captity))
	case leakBucketLimter:
		logger.Info("leakBucketLimter")
		limiter = NewLeakBucketLimiter(lConfig.Name, lConfig.Captity, lConfig.Duration)
	case slideWindowLimiter:
		logger.Info("slideWindowLimiter not has")
	default:
		str := fmt.Sprintf("Error Type %d", t)
		logger.Error(str)
	}

	return limiter
}

package proxy

type LimitType int

const (
	queueLimter  LimitType = iota + 1 // 队列模式限流
	bucketLimter                      // golang官方库  golang.org/x/time/rate  bucket
	slideWindow                       // 滑动窗口限流器 tcp滑动窗口
)

// 限流接口
type LimitInterface interface {
	IsAvalivale() bool
	Run(handler func(conn interface{}))
	SetWaitQueue(conn interface{})
}

///////////////////////////////QueueLimit/////////////////////////////
// 通过队列实现限流
type QueueLimit struct {
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
func NewQueueLimter(waitLength, maxConn int) *QueueLimit {
	limiter := &QueueLimit{
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
func (this QueueLimit) IsAvalivale() bool {
	length := len(this.waitQueue)

	// 超过了等待队列的数量 说明超过了最大并发持续进行中
	if length >= this.initWaitQueueLen {
		return false
	}

	return true
}

// 限流器增加计数
func (this *QueueLimit) SetWaitQueue(conn interface{}) {
	this.waitQueue <- conn
}

func (this QueueLimit) Run(handler func(conn interface{})) {
	go func() {
		for connection := range this.waitQueue {
			<-this.availPools
			go func(connection interface{}) {
				handler(connection)
				this.availPools <- true
				logger.Info("conn handle ok on QueueLimiter")
			}(connection)
		}
	}()

}

///////////////////////////////QueueLimit/////////////////////////////

package base

import (
	"sync"
	"sync/atomic"
)

// 单例
type singleton struct {
}

var instance *singleton

// 第一种方法-线程安全
var mu sync.Mutex

func GetInstance() *singleton {
	if instance == nil { // <-- Not yet perfect. since it's not fully atomic
		mu.Lock()
		defer mu.Unlock()

		if instance == nil {
			instance = &singleton{}
		}
	}
	return instance
}

// 第二种线程安全
var initialized uint32

func GetInstanceAtomic() *singleton {
	if atomic.LoadUint32(&initialized) == 1 {
		return instance
	}

	mu.Lock()
	defer mu.Unlock()

	if initialized == 0 {
		instance = &singleton{}
		atomic.StoreUint32(&initialized, 1)
	}

	return instance
}

// 第三种线程安全的单利
var once sync.Once

func GetInstanceOnce() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}

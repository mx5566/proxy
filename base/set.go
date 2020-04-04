package base

import "sync"

// set 类似c++ STL
type ISet interface {
}

// Set
type Set struct {
	m map[ISet]bool
	sync.RWMutex
}

func New() *Set {
	return &Set{m: make(map[ISet]bool)}
}

// Add
func (this *Set) Add(s ISet) {
	this.Lock()
	this.m[s] = true
	this.Unlock()
}

// Test
func (this *Set) Test(s ISet) bool {
	if _, ok := this.m[s]; ok {
		return true
	}
	return false
}

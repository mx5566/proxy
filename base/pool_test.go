// Copyright (c) 2020 by meng.  All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

/**
 * @Author: meng
 * @Description:
 * @File:  pool_test
 * @Version: 1.0.0
 * @Date: 2020/4/4 21:40
 */

package base

import (
	"fmt"
	"testing"
)

type JobEat struct {
}

func (this *JobEat) Do() error {
	fmt.Println("JonEat Do.......")
	return nil
}

func TestNewWorker(t *testing.T) {
	dispatcher := NewDispatcher(10)
	dispatcher.Run()

	job := &JobEat{}
	JobQueue <- job
}

func BenchmarkNewWorker(b *testing.B) {
	dispatcher := NewDispatcher(10)
	dispatcher.Run()

	job := &JobEat{}

	JobQueue <- job
}

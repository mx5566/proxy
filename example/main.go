// Copyright 2020 mengxiang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

/**
 * @Author: meng
 * @Description:
 * @File:  main.go
 * @Version: 1.0.0
 * @Date: 2020/4/5 15:46
 */

package main

import (
	"fmt"

	"github.com/mx5566/proxy"
)

func main() {
	fmt.Print("HelloWorld")
	proxy.CreateProxy("../config.yaml")
}

// Copyright (c) 2020 by meng.  All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"testing"

	"fmt"
)

func TestNew(t *testing.T) {
	set := New()

	set.Add(1)
	set.Add(1)
	set.Add(2)

	fmt.Println(set.m)

	fmt.Println(set.Test(1))
}

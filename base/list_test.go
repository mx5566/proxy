// Copyright (c) 2020 by meng.  All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"fmt"
	"testing"
)

func TestNewLinkList(t *testing.T) {
	linkList := NewLinkList()

	node := new(Node)
	node.next = nil
	node.prev = nil
	node.value = 1

	linkList.Add(node)

	node = new(Node)
	node.value = 2
	linkList.Add(node)

	head := linkList.head
	for head != nil {
		fmt.Print(head.value)
		head = head.next
	}
}

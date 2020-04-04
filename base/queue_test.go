package base

import (
	"fmt"
	"testing"
)

func TestNewLinkQueue(t *testing.T) {
	linkQueue := NewLinkQueue()

	node := new(Node)
	node.next = nil
	node.prev = nil
	node.value = 1

	linkQueue.Push(node)

	node = new(Node)
	node.value = 2
	linkQueue.Push(node)

	fmt.Println(linkQueue.IsEmpty())

	nodeRet := linkQueue.Pop()

	fmt.Println(nodeRet.value)
}

func TestNewCircleQueue(t *testing.T) {
	circleQueue := NewCircleQueue(3)
	circleQueue.Push(1)
	circleQueue.Push(2)
	circleQueue.Push(3)
	circleQueue.Push(4)

	fmt.Println(circleQueue.IsFull())

	fmt.Println(circleQueue.Pop())
	fmt.Println(circleQueue.Pop())
	fmt.Println(circleQueue.Pop())
	fmt.Println(circleQueue.Pop())
	fmt.Println(circleQueue.Pop())

	fmt.Println(circleQueue.IsEmpty())
}

func BenchmarkNewCircleQueue(b *testing.B) {
	circleQueue := NewCircleQueue(3)
	circleQueue.Push(1)
	circleQueue.Push(2)
	circleQueue.Push(3)
	circleQueue.Push(4)

	fmt.Println(circleQueue.IsFull())

	fmt.Println(circleQueue.Pop())
	fmt.Println(circleQueue.Pop())
	fmt.Println(circleQueue.Pop())
	fmt.Println(circleQueue.Pop())
	fmt.Println(circleQueue.Pop())

	fmt.Println(circleQueue.IsEmpty())
}

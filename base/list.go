package base

// 节点值
type NodeValue interface {
}

// 链表节点
type Node struct {
	// 节点值
	value NodeValue
	// 前向节点
	prev *Node
	// 后一个节点
	next *Node
}

type LinkList struct {
	// 链表头
	head *Node
	// 链表尾
	tail *Node
}

func NewLinkList() *LinkList {
	// 初始化链表
	ret := &LinkList{
		head: nil,
		tail: nil,
	}

	return ret
}

// 增加节点 从尾插入
func (this *LinkList) Add(node *Node) {
	if node == nil {
		return
	}

	// 增加节点 空链表
	if this.tail == nil {
		this.head = node
		this.tail = node

		return
	}

	this.tail.next = node
	node.prev = this.tail

	// 尾节点指向插入节点
	this.tail = node
}

// 删除节点头结点
func (this *LinkList) Delete() {
	if this.head == nil {
		return
	}

	// 只有一个节点
	if this.head.next == nil {
		this.head, this.tail = nil, nil
		return
	}
	//
	temp := this.head.next

	this.head.next.prev = nil
	this.head.next = nil

	this.head = temp

}

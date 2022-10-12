package hw04lrucache

type List interface {
	Len() int
	Front() *Node
	Back() *Node
	PushFront(v interface{}) *Node
	PushBack(v interface{}) *Node
	Remove(i *Node)
	MoveToFront(i *Node)
}

type Node struct {
	Next  *Node
	Value any
	Prev  *Node
}

type DubleList struct {
	Head *Node
	Tail *Node
	len  int
}

func (dl DubleList) Len() int {
	return dl.len
}

func (dl DubleList) Front() *Node {
	return dl.Head
}

func (dl DubleList) Back() *Node {
	return dl.Tail
}

func (dl *DubleList) PushFront(newValue interface{}) *Node {
	newData := &Node{Value: newValue}
	if dl.Head == nil {
		dl.Head = newData
		dl.Tail = newData
	} else {
		newData.Next = dl.Head
		dl.Head.Prev = newData
		dl.Head = newData
	}
	dl.len++
	return newData
}

func (dl *DubleList) Remove(item *Node) {
	if item.Prev == nil {
		dl.Head = dl.Head.Next
		dl.Head.Prev = nil
	} else if item.Next == nil {
		dl.Tail = dl.Tail.Prev
		dl.Tail.Next = nil
	} else {
		item.Next.Prev = item.Prev
		item.Prev.Next = item.Next
	}
	dl.len--
}

func (dl *DubleList) PushBack(newValue interface{}) *Node {
	newData := &Node{Value: newValue}
	if dl.Head == nil {
		dl.Head = newData
		dl.Tail = newData
	} else {
		dl.Tail.Next = newData
		newData.Prev = dl.Tail
		dl.Tail = newData
	}
	dl.len++
	return newData
}

func (dl *DubleList) MoveToFront(item *Node) {
	dl.Remove(item)
	oldHead := dl.Head
	dl.Head.Prev = item
	dl.Head = item
	dl.Head.Prev = nil
	dl.Head.Next = oldHead
	dl.len++
}

func NewList() List {
	return List(&DubleList{})
}

package linkedlist

import "fmt"

type Node struct {
	Value int
	Next  *Node
}

type LinkedList struct {
	Head *Node
}

type AccessList interface {
	InsertNode(int)
	ReadList()
	ReverseList()
}

func (l *LinkedList) InsertNode(value int) {
	newNode := Node{
		Value: value,
		Next:  nil,
	}

	if l.Head == nil {
		l.Head = &newNode
		return
	}

	temp := l.Head
	for temp.Next != nil {
		temp = temp.Next
	}
	temp.Next = &newNode
}

func (l *LinkedList) ReadList() {
	if l.Head == nil {
		return
	}

	temp := l.Head
	for temp != nil {
		fmt.Println(temp.Value)
		temp = temp.Next
	}
}

func (l *LinkedList) ReverseList() {
	if l.Head == nil {
		return
	}

	var last *Node
	current := l.Head
	for current != nil {
		next := current.Next
		current.Next = last
		last = current
		current = next
	}
	l.Head = last
}

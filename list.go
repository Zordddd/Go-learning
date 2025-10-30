package main

import (
	"fmt"
)


type Node struct {
    value int
    next *Node
}

type List struct {
    head *Node
    tail *Node
}

func (l *List) AddNode(value int) {
    newNode := &Node{value : value}

    if l.head == nil {
        l.head = newNode
        l.tail = newNode
    } else {
        l.tail.next = newNode
        l.tail = newNode
    }
}

func (l *List) Print() {
    currentNode := l.head
    for currentNode != nil {
        fmt.Printf("%d ", currentNode.value)
        currentNode = currentNode.next
    }
    fmt.Printf("\n")
}

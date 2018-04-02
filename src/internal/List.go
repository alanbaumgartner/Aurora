package internal

import (
	"net"
)

type List struct {
	List []net.Conn
}

func NewList() List {
	self := List{}
	self.List = []net.Conn{}
	return self
}

func (list *List) Get(index int) net.Conn {
	if index < len(list.List) {
		return list.List[index]
	}
	return nil
}

func (list *List) Add(conn net.Conn) {
	for _, conn := range list.List {
		if conn == conn {
			return
		}
	}
	list.List = append(list.List, conn)
}

func (list *List) Remove(conn net.Conn) {
	for index, conn := range list.List {
		if conn == conn {
			list.List = append(list.List[:index], list.List[index+1:]...)
			return
		}
	}
}

func (list *List) All() []net.Conn {
	return list.List
}

func (list *List) Clear() {
	list.List = []net.Conn{}
}

func (list *List) isEmpty() bool {
	if len(list.List) == 0 {
		return true
	}
	return false
}

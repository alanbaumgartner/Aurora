package internal

import (
	"net"
	"sync"
)

type Conn net.Conn

type Set struct {
	items map[Conn]bool
	lock  sync.RWMutex
}

func (set *Set) Add(conn Conn) *Set {
	set.lock.Lock()
	defer set.lock.Unlock()
	if set.items == nil {
		set.items = make(map[Conn]bool)
	}
	_, ok := set.items[conn]
	if !ok {
		set.items[conn] = true
	}
	return set
}

func (set *Set) Clear() {
	set.lock.Lock()
	defer set.lock.Unlock()
	set.items = make(map[Conn]bool)
}

func (set *Set) Delete(conn Conn) bool {
	set.lock.Lock()
	defer set.lock.Unlock()
	_, ok := set.items[conn]
	if ok {
		delete(set.items, conn)
	}
	return ok
}

func (set *Set) Has(conn Conn) bool {
	set.lock.RLock()
	defer set.lock.RUnlock()
	_, ok := set.items[conn]
	return ok
}

func (set *Set) Items() []Conn {
	set.lock.RLock()
	defer set.lock.RUnlock()
	var items []Conn
	for i := range set.items {
		items = append(items, i)
	}
	return items
}

func (set *Set) Size() int {
	set.lock.RLock()
	defer set.lock.RUnlock()
	return len(set.items)
}

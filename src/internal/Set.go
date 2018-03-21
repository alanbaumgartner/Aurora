package internal

import (
	"net"
	"sync"
)

type Conn net.Conn

type ItemSet struct {
	items map[Conn]bool
	lock  sync.RWMutex
}

func (s *ItemSet) Add(conn Conn) *ItemSet {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.items == nil {
		s.items = make(map[Conn]bool)
	}
	_, ok := s.items[conn]
	if !ok {
		s.items[conn] = true
	}
	return s
}

func (s *ItemSet) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = make(map[Conn]bool)
}

func (s *ItemSet) Delete(conn Conn) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.items[conn]
	if ok {
		delete(s.items, conn)
	}
	return ok
}

func (s *ItemSet) Has(conn Conn) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.items[conn]
	return ok
}

func (s *ItemSet) Items() []Conn {
	s.lock.RLock()
	defer s.lock.RUnlock()
	var items []Conn
	for i := range s.items {
		items = append(items, i)
	}
	return items
}

func (s *ItemSet) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}

package store

import (
	"sync"
)

var DefaultStore Store

// Store - simple key-val database
type Store struct {
	sync.RWMutex
	mp map[string][]byte
}

func init() {
	DefaultStore = New()
}

func New() Store {
	return Store{
		mp: make(map[string][]byte),
	}
}

func (s *Store) Set(key string, val []byte) {
	s.Lock()
	defer s.Unlock()
	s.mp[key] = val
}

func (s *Store) Get(key string) []byte {
	s.RLock()
	defer s.RUnlock()
	return s.mp[key]
}

package store

import (
	"sync"
)

type Store struct {
	dict map[string]int64
	mu   sync.Mutex
}

func NewStore() *Store {
	return &Store{
		dict: make(map[string]int64),
	}
}
func (store *Store) Add(key string, value int64) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.dict[key] = value
}
func (store *Store) Set(key string) (int64, bool) {
	store.mu.Lock()
	defer store.mu.Unlock()

	value, exist := store.dict[key]

	return value, exist
}

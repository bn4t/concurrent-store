package concurrent_store

import (
	"sync"
)

// Store is a concurrency-safe deduplicated store
type Store struct {
	lock  *sync.RWMutex
	items map[interface{}]struct{}
}

func NewStore() *Store {
	return &Store{
		lock:  &sync.RWMutex{},
		items: map[interface{}]struct{}{},
	}
}

func (l *Store) Add(v interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.items[v] = struct{}{}
}

func (l *Store) Contains(v interface{}) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()

	_, ok := l.items[v]

	return ok
}

func (l *Store) Pop() (interface{}, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	for k := range l.items {
		delete(l.items, k)
		return k, nil
	}

	return nil, ErrStoreEmpty
}

func (l *Store) All() map[interface{}]struct{} {
	l.lock.RLock()
	defer l.lock.RUnlock()

	// clone the map
	ret := map[interface{}]struct{}{}
	for k, v := range l.items {
		ret[k] = v
	}
	return ret
}

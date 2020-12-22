package concurrent_store

import (
	"sync"
)

// Map is a concurrency-safe map implementation
type Map struct {
	lock  *sync.RWMutex
	items map[interface{}]interface{}
}

func NewMap() *Map {
	return &Map{
		lock:  &sync.RWMutex{},
		items: map[interface{}]interface{}{},
	}
}

func (l *Map) Add(k interface{}, v interface{}) {
	l.lock.Lock()
	l.items[k] = v
	l.lock.Unlock()
}

func (l *Map) Contains(k interface{}) bool {
	l.lock.RLock()
	_, ok := l.items[k]
	l.lock.RUnlock()

	return ok
}

func (l *Map) Pop() (interface{}, interface{}, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if len(l.items) == 0 {
		return nil, nil, ErrStoreEmpty
	}

	for k, v := range l.items {
		delete(l.items, k)
		return k, v, nil
	}

	panic("reached unreachable state in concurrent map")
	return nil, nil, nil
}

func (l *Map) All() map[interface{}]interface{} {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.items
}

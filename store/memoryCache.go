package store

import (
	"errors"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"sync"
)

func RlockAndUnlock(lock *sync.RWMutex) func() {
	lock.RLock()
	return func() {
		lock.RUnlock()
	}
}

func lockAndUnlock(lock *sync.RWMutex) func() {
	lock.Lock()
	return func() {
		lock.Unlock()
	}
}

type Memory struct {
	cache map[string]*linkedhashmap.Map
	lock  *sync.RWMutex
}

func NewMemory() *Memory {
	return &Memory{
		cache: make(map[string]*linkedhashmap.Map),
		lock:  &sync.RWMutex{},
	}
}

func (m *Memory) getCache(code string) (cache *linkedhashmap.Map) {
	if _, ok := m.cache[code]; !ok {
		m.cache[code] = linkedhashmap.New()
	}
	return m.cache[code]
}

func (m *Memory) GetLast(code string) (string, interface{}) {
	defer RlockAndUnlock(m.lock)()
	cache := m.getCache(code)
	iterator := cache.Iterator()
	if !iterator.Last() {
		return "", nil
	}
	return iterator.Key().(string), iterator.Value()
}

func (m *Memory) GetString(code, key string) (string, error) {
	v, err := m.Get(code, key)
	if err != nil {
		return "", err
	}
	if vv, ok := v.(string); !ok {
		return vv, nil
	}
	return "", errors.New("not string")
}

func (m *Memory) Get(code, key string) (interface{}, error) {
	defer RlockAndUnlock(m.lock)()
	cache := m.getCache(code)
	v, ok := cache.Get(key)
	if !ok {
		return "", errors.New("not found")
	}
	return v.(string), nil
}

func (m *Memory) Set(code, key string, value interface{}) error {
	defer lockAndUnlock(m.lock)()
	cache := m.getCache(code)
	cache.Put(key, value)
	return nil
}

func (m *Memory) Delete(code, key string) {
	defer lockAndUnlock(m.lock)()
	cache := m.getCache(code)
	cache.Remove(key)
}

package server

import (
	"sync"
	"reflect"
)


type Memory struct {
	mu    sync.RWMutex
	items map[string]interface{}
}

func NewMemory() Memory {
	return Memory{items: make(map[string]interface{})}
}

func (m *Memory) Get(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, found := m.items[key]
	if ! found {
		return nil, false
	}

	return item, true
}

func (m *Memory) GetW(key string, defaultValue ...interface{}) interface{} {
	item, found := m.Get(key)
	if found {
		return item
	}

	if len(defaultValue) < 1 {
		return nil
	}

	return defaultValue[0]
}

func (m *Memory) GetDelete(key string) (interface{}, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, found := m.Get(key)
	if found {
		m.Delete(key)
	}

	return item, found
}

func (m *Memory) Put(key string, value interface{}) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[key] = value

	return true
}

func (m *Memory) Update(key string, value interface{}) bool {
	return m.Put(key, value)
}

func (m *Memory) Delete(k string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, found := m.items[k]
	if found {
		delete(m.items, k)
	}

	return found
}

func (m *Memory) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = make(map[string]interface{})
}


/**
 * Array helpers
 */
func (m *Memory) PutToArray(key string, value ...interface{}) bool {
	items, found := m.GetDelete(key)
	if !found {
		items = make([]interface{}, 0)
	}

	itemType := reflect.TypeOf(items).Kind()
	if itemType != reflect.Array && itemType != reflect.Slice {
		return false
	}

	items = append(items.([]interface{}), value...)

	return m.Put(key, items)
}
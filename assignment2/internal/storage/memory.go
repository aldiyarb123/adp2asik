package storage

import "sync"

type MemoryStore struct {
	mu   sync.Mutex
	data map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]string),
	}
}

func (m *MemoryStore) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *MemoryStore) GetAll() map[string]string {
	m.mu.Lock()
	defer m.mu.Unlock()

	copy := make(map[string]string)
	for k, v := range m.data {
		copy[k] = v
	}
	return copy
}

func (m *MemoryStore) Delete(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[key]; !ok {
		return false
	}
	delete(m.data, key)
	return true
}

func (m *MemoryStore) Size() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.data)
}

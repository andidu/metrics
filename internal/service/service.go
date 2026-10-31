package service

import (
	"database/sql"
	"sync"

	"github.com/andidu/metrics/internal/handler"
)

func NewMemStorage(db *sql.DB) handler.MemStorage {
	return &memStorageImpl{
		counters: make(map[string]int64),
		gauges:   make(map[string]float64),
		db:       db,
	}
}

type memStorageImpl struct {
	m        sync.Mutex
	counters map[string]int64
	gauges   map[string]float64
	db       *sql.DB
}

func (m *memStorageImpl) Gauges() map[string]float64 {
	m.m.Lock()
	defer m.m.Unlock()

	gauges := make(map[string]float64, len(m.gauges))
	for k, v := range m.gauges {
		gauges[k] = v
	}
	return gauges
}

func (m *memStorageImpl) Counters() map[string]int64 {
	m.m.Lock()
	defer m.m.Unlock()

	counters := make(map[string]int64, len(m.counters))
	for k, v := range m.counters {
		counters[k] = v
	}
	return counters
}

func (m *memStorageImpl) UpdateGauge(name string, value float64) error {
	m.m.Lock()
	m.gauges[name] = value
	m.m.Unlock()
	return nil
}

func (m *memStorageImpl) UpdateCounter(name string, value int64) error {
	m.m.Lock()
	m.counters[name] += value
	m.m.Unlock()
	return nil
}

func (m *memStorageImpl) OverrideCounter(name string, value int64) error {
	m.m.Lock()
	m.counters[name] = value
	m.m.Unlock()
	return nil
}

func (m *memStorageImpl) Ping() bool {
	return m.db.Ping() == nil
}

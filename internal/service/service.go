package service

import (
	"github.com/andidu/metrics/internal/handler"
)

func NewMemStorage() handler.MemStorage {
	return memStorageImpl{
		counters: make(map[string]int64),
		gauges:   make(map[string]float64),
	}
}

type memStorageImpl struct {
	counters map[string]int64
	gauges   map[string]float64
}

func (m memStorageImpl) Gauges() map[string]float64 {
	return m.gauges
}

func (m memStorageImpl) Counters() map[string]int64 {
	return m.counters
}

func (m memStorageImpl) UpdateGauge(name string, value float64) error {
	m.gauges[name] = value
	return nil
}

func (m memStorageImpl) UpdateCounter(name string, value int) error {
	m.counters[name] += int64(value)
	return nil
}

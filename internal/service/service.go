package service

import (
	"strconv"
	"strings"
)

type MemStorage interface {
	Gauges() map[string]float64
	Counters() map[string]int
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int)
	GetGauge(name string) (string, bool)
	GetCounter(name string) (string, bool)
}

func NewMemStorage() MemStorage {
	return memStorageImpl{
		counters: make(map[string]int),
		gauges:   make(map[string]float64),
	}
}

type memStorageImpl struct {
	counters map[string]int
	gauges   map[string]float64
}

func (m memStorageImpl) Gauges() map[string]float64 {
	return m.gauges
}

func (m memStorageImpl) Counters() map[string]int {
	return m.counters
}

func (m memStorageImpl) UpdateGauge(name string, value float64) {
	m.gauges[name] = value
}

func (m memStorageImpl) UpdateCounter(name string, value int) {
	m.counters[name] += value
}

func (m memStorageImpl) GetGauge(name string) (string, bool) {
	fval, ok := m.gauges[name]
	if ok {
		s := strconv.FormatFloat(fval, 'f', 3, 64)
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
		return s, true
	}
	return "", false
}

func (m memStorageImpl) GetCounter(name string) (string, bool) {
	ival, ok := m.counters[name]
	if ok {
		return strconv.Itoa(ival), true
	}
	return "", false
}

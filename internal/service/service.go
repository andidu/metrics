package service

import (
	"errors"
	"strconv"
	"strings"

	"github.com/andidu/metrics/internal/handler"
)

var noElementErr = errors.New("No such element found")

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

func (m memStorageImpl) GetGauge(name string) (string, error) {
	fval, ok := m.gauges[name]
	if ok {
		s := strconv.FormatFloat(fval, 'f', 3, 64)
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
		return s, nil
	}
	return "", noElementErr
}

func (m memStorageImpl) GetCounter(name string) (string, error) {
	ival, ok := m.counters[name]
	if ok {
		return strconv.FormatInt(ival, 10), nil
	}
	return "", noElementErr
}

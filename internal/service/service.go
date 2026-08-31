package service

import (
	"strconv"
)

type MemStorage struct {
	Counters map[string]int
	Gauges   map[string]float64
}

var TmpInMemoryStarage = MemStorage{
	Counters: make(map[string]int),
	Gauges:   make(map[string]float64),
}

func (MemStorage) UpdateGauge(name string, value float64) {
	TmpInMemoryStarage.Gauges[name] = value
}

func (MemStorage) UpdateCounter(name string, value int) {
	TmpInMemoryStarage.Counters[name] += value
}

func (MemStorage) GetGauge(name string) (string, bool) {
	fval, ok := TmpInMemoryStarage.Gauges[name]
	if ok {
		return strconv.FormatFloat(fval, 'f', 2, 64), true
	}
	return "", false
}

func (MemStorage) GetCounter(name string) (string, bool) {
	ival, ok := TmpInMemoryStarage.Counters[name]
	if ok {
		return strconv.Itoa(ival), true
	}
	return "", false
}

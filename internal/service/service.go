package service

type MemStorage struct {
	counters map[string]int
	gauges   map[string]float64
}

var TmpInMemoryStarage = MemStorage{
	counters: make(map[string]int),
	gauges:   make(map[string]float64),
}

func (MemStorage) UpdateGauge(name string, value float64) {
	TmpInMemoryStarage.gauges[name] = value
}

func (MemStorage) UpdateCounter(name string, value int) {
	TmpInMemoryStarage.counters[name] += value
}

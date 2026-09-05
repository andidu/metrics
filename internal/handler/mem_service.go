package handler

type MemStorage interface {
	Gauges() map[string]float64
	Counters() map[string]int64
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int)
	GetGauge(name string) (string, bool)
	GetCounter(name string) (string, bool)
}

package handler

type MemStorage interface {
	Gauges() map[string]float64
	Counters() map[string]int64
	UpdateGauge(name string, value float64) error
	UpdateCounter(name string, value int) error
	OverrideCounter(name string, value int64) error
}

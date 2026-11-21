package agent

import (
	"maps"
	"math/rand"
	"runtime"
)

type MetricsSample struct {
	Counters map[string]int64
	Gauges   map[string]float64
}

func (m MetricsSample) Clone() MetricsSample {
	return MetricsSample{
		Counters: maps.Clone(m.Counters),
		Gauges:   maps.Clone(m.Gauges),
	}
}

func (m MetricsSample) AddGauges(g map[string]float64) MetricsSample {
	if m.Gauges == nil {
		m.Gauges = make(map[string]float64)
	}

	maps.Copy(m.Gauges, g)
	return m
}

func (m MetricsSample) AddPollCount(counter int64) MetricsSample {
	if m.Counters == nil {
		m.Counters = make(map[string]int64)
	}

	m.Counters["PollCount"] = counter
	return m
}

func (m MetricsSample) PollCount() int64 {
	return m.Counters["PollCount"]
}

func ObtainMetricsSample() MetricsSample {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	return MetricsSample{
		Counters: map[string]int64{},
		Gauges: map[string]float64{
			"Alloc":         float64(stats.Alloc),
			"BuckHashSys":   float64(stats.BuckHashSys),
			"Frees":         float64(stats.Frees),
			"GCCPUFraction": stats.GCCPUFraction,
			"GCSys":         float64(stats.GCSys),
			"HeapAlloc":     float64(stats.HeapAlloc),
			"HeapIdle":      float64(stats.HeapIdle),
			"HeapInuse":     float64(stats.HeapInuse),
			"HeapObjects":   float64(stats.HeapObjects),
			"HeapReleased":  float64(stats.HeapReleased),
			"HeapSys":       float64(stats.HeapSys),
			"LastGC":        float64(stats.LastGC),
			"Lookups":       float64(stats.Lookups),
			"MCacheInuse":   float64(stats.MCacheInuse),
			"MCacheSys":     float64(stats.MCacheSys),
			"MSpanInuse":    float64(stats.MSpanInuse),
			"MSpanSys":      float64(stats.MSpanSys),
			"Mallocs":       float64(stats.Mallocs),
			"NextGC":        float64(stats.NextGC),
			"NumForcedGC":   float64(stats.NumForcedGC),
			"NumGC":         float64(stats.NumGC),
			"OtherSys":      float64(stats.OtherSys),
			"PauseTotalNs":  float64(stats.PauseTotalNs),
			"StackInuse":    float64(stats.StackInuse),
			"StackSys":      float64(stats.StackSys),
			"Sys":           float64(stats.Sys),
			"TotalAlloc":    float64(stats.TotalAlloc),
			"RandomValue":   rand.Float64(),
		},
	}
}

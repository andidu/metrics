package agent_test

import (
	"testing"

	"github.com/andidu/metrics/internal/agent"
	"github.com/stretchr/testify/assert"
)

func TestObtainMetricsSample(t *testing.T) {
	tests := []struct {
		name     string // description of this test case
		counters []string
		gauges   []string
	}{
		{
			name: "check fields",
			counters: []string{
				"counter",
			},
			gauges: []string{
				"Alloc",
				"BuckHashSys",
				"Frees",
				"GCCPUFraction",
				"GCSys",
				"HeapAlloc",
				"HeapIdle",
				"HeapInuse",
				"HeapObjects",
				"HeapReleased",
				"HeapSys",
				"LastGC",
				"Lookups",
				"MCacheInuse",
				"MCacheSys",
				"MSpanInuse",
				"MSpanSys",
				"Mallocs",
				"NextGC",
				"NumForcedGC",
				"NumGC",
				"OtherSys",
				"PauseTotalNs",
				"StackInuse",
				"StackSys",
				"Sys",
				"TotalAlloc",
				"RandomValue",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := agent.ObtainMetricsSample()

			countersKeys := make([]string, 0, len(got.Counters))
			for k := range got.Counters {
				countersKeys = append(countersKeys, k)
			}

			gaugeKeys := make([]string, 0, len(got.Gauges))
			for k := range got.Gauges {
				gaugeKeys = append(gaugeKeys, k)
			}

			assert.ElementsMatch(t, countersKeys, tt.counters)
			assert.ElementsMatch(t, gaugeKeys, tt.gauges)
		})
	}
}

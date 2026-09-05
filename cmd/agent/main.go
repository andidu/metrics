package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/andidu/metrics/internal/agent"
)

func main() {
	var flags = parseFlags()

	var gaugeUrlTemplate = fmt.Sprintf("http://%s/update/gauge", flags.serverAddress) + "/%s/%f"
	var counterUrlTemplate = fmt.Sprintf("http://%s/update/counter", flags.serverAddress) + "/%s/%d"

	var mutex sync.Mutex
	var sample agent.MetricsSample

	go func() {
		var counter = int64(0)
		for {
			metrics := agent.ObtainMetricsSample(counter)
			mutex.Lock()
			sample = metrics
			mutex.Unlock()
			time.Sleep(time.Duration(flags.metrics.pollInterval) * time.Second)

			counter += 1
		}
	}()

	for {
		time.Sleep(time.Duration(flags.metrics.repeatInterval) * time.Second)
		mutex.Lock()
		metrics := sample
		mutex.Unlock()
		for name, value := range metrics.Counters {
			_, err := http.Post(fmt.Sprintf(counterUrlTemplate, name, value), "text/plain", nil)
			if err != nil {
				fmt.Println(err.Error())
			}
		}

		for name, value := range metrics.Gauges {
			_, err := http.Post(fmt.Sprintf(gaugeUrlTemplate, name, value), "text/plain", nil)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

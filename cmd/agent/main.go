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

	var mutex sync.Mutex // user for guarding sample, counter and sentCounter
	var sample agent.MetricsSample
	counter := int64(0)

	go func() {
		for {
			var localCounter int64
			mutex.Lock()
			localCounter = counter
			mutex.Unlock()

			metrics := agent.ObtainMetricsSample(localCounter)

			mutex.Lock()
			sample = metrics
			counter += 1
			mutex.Unlock()
			time.Sleep(time.Duration(flags.metrics.pollInterval) * time.Second)
		}
	}()

	for {
		time.Sleep(time.Duration(flags.metrics.repeatInterval) * time.Second)
		mutex.Lock()
		metrics := sample
		mutex.Unlock()
		for name, value := range metrics.Counters {
			resp, err := http.Post(fmt.Sprintf(counterUrlTemplate, name, value), "text/plain", nil)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				defer resp.Body.Close()
				mutex.Lock()
				counter -= value
				metrics.InvalidateCounter(name)
				mutex.Unlock()
			}
		}

		for name, value := range metrics.Gauges {
			resp, err := http.Post(fmt.Sprintf(gaugeUrlTemplate, name, value), "text/plain", nil)
			if err != nil {
				fmt.Println(err.Error())
			}
			defer resp.Body.Close()
		}
	}
}

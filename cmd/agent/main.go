package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/andidu/metrics/internal/agent"
)

func main() {
	flag.Parse()

	var gaugeUrlTemplate = fmt.Sprintf("http://%s/update/gauge", *serverAddress) + "/%s/%f"
	var counterUrlTemplate = fmt.Sprintf("http://%s/update/counter", *serverAddress) + "/%s/%d"

	var mutex sync.Mutex
	var sample agent.MetricsSample

	go func() {
		for true {
			metrics := agent.ObtainMetricsSample()
			mutex.Lock()
			sample = metrics
			mutex.Unlock()
			time.Sleep(time.Duration(*pollInterval) * time.Second)
		}
	}()

	for true {
		time.Sleep(time.Duration(*repeatInterval) * time.Second)
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

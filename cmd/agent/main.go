package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/andidu/metrics/internal/agent"
	config "github.com/andidu/metrics/internal/config/agent"
)

func main() {
	var flags = config.ParseConfig()

	var gaugeUrlTemplate = fmt.Sprintf("http://%s/update/gauge", flags.ServerAddress) + "/%s/%f"
	var counterUrlTemplate = fmt.Sprintf("http://%s/update/counter", flags.ServerAddress) + "/%s/%d"

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
			time.Sleep(time.Duration(flags.Metrics.PollInterval) * time.Second)
		}
	}()

	for {
		time.Sleep(time.Duration(flags.Metrics.RepeatInterval) * time.Second)
		mutex.Lock()
		metrics := sample
		mutex.Unlock()
		for name, value := range metrics.Counters {
			resp, err := http.Post(fmt.Sprintf(counterUrlTemplate, name, value), "text/plain", nil)
			if err != nil {
				log.Println(err.Error())
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
				log.Println(err.Error())
			}
			defer resp.Body.Close()
		}
	}
}

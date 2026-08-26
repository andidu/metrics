package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andidu/metrics/internal/agent"
)

func main() {

	var gaugeUrlTemplate = "http://localhost:8080/update/gauge/%s/%d"
	var counterUrlTemplate = "http://localhost:8080/update/counter/%s/%f"
	var counter = 0
	for true {
		var metrics = agent.ObtainMetricsSample()

		if counter%5 == 0 {
			for name, value := range metrics.Counters {
				http.Post(fmt.Sprintf(counterUrlTemplate, name, value), "text/plain", nil)
			}

			for name, value := range metrics.Gauges {
				http.Post(fmt.Sprintf(gaugeUrlTemplate, name, value), "text/plain", nil)
			}
		}

		counter += 1
		time.Sleep(2 * time.Second)
	}
}

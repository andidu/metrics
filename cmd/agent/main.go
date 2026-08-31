package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andidu/metrics/internal/agent"
)

func main() {

	var gaugeUrlTemplate = "http://localhost:8080/update/gauge/%s/%f"
	var counterUrlTemplate = "http://localhost:8080/update/counter/%s/%d"
	var counter = 0
	for true {
		var metrics = agent.ObtainMetricsSample()

		if counter%5 == 0 {
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

		counter += 1
		time.Sleep(2 * time.Second)
	}
}

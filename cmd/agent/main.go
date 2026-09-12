package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/andidu/metrics/internal/agent"
	config "github.com/andidu/metrics/internal/config/agent"
	models "github.com/andidu/metrics/internal/model"
)

func main() {
	var flags = config.ParseConfig()

	var URLTemplate = fmt.Sprintf("http://%s/update/", flags.ServerAddress)

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
			v := int64(value)
			var m = models.Metrics{
				ID:    name,
				MType: models.Counter,
				Delta: &v,
			}
			body, err := json.Marshal(m)
			if err != nil {
				log.Println(err.Error())
				return
			}

			resp, err := http.Post(URLTemplate, "application/json", bytes.NewReader(body))
			if err != nil {
				log.Println(err.Error())
			} else {
				mutex.Lock()
				counter -= value
				metrics.InvalidateCounter(name)
				mutex.Unlock()
				resp.Body.Close()
			}
		}

		for name, value := range metrics.Gauges {
			var m = models.Metrics{
				ID:    name,
				MType: models.Gauge,
				Value: &value,
			}
			body, err := json.Marshal(m)
			if err != nil {
				log.Println(err.Error())
				return
			}

			resp, err := http.Post(URLTemplate, "application/json", bytes.NewReader(body))
			if err != nil {
				log.Println(err.Error())
			}
			if err == nil {
				resp.Body.Close()
			}
		}
	}
}

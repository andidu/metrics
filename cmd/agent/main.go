package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/andidu/metrics/internal/agent"
	config "github.com/andidu/metrics/internal/config/agent"
	models "github.com/andidu/metrics/internal/model"
	"github.com/andidu/metrics/internal/utils"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"golang.org/x/sync/semaphore"
)

func main() {
	var flags, err = config.ParseConfig()
	if err != nil {
		log.Println(err.Error())
		return
	}

	var URLTemplate = fmt.Sprintf("http://%s/updates/", flags.ServerAddress)

	var mutex sync.Mutex // user for guarding sample, counter and sentCounter
	var sample agent.MetricsSample
	counter := int64(0)

	// Extra gauge values
	var extraGauges = make(map[string]float64)

	var errorHappened = false
	var pullSignalExtra = make(chan struct{})
	var pullSignalBasic = make(chan struct{})

	go func() {
		for {
			<-pullSignalExtra

			if errorHappened {
				return
			}

			stats, err := mem.VirtualMemory()
			if err != nil {
				mutex.Lock()
				for k := range extraGauges {
					delete(extraGauges, k)
				}
				mutex.Unlock()

				continue
			}

			percentages, err := cpu.Percent(0, true)
			if err != nil {
				mutex.Lock()
				for k := range extraGauges {
					delete(extraGauges, k)
				}
				extraGauges["TotalMemory"] = float64(stats.Total)
				extraGauges["FreeMemory"] = float64(stats.Free)
				mutex.Unlock()

				continue
			}

			mutex.Lock()
			extraGauges["TotalMemory"] = float64(stats.Total)
			extraGauges["FreeMemory"] = float64(stats.Free)
			for i, x := range percentages {
				extraGauges[fmt.Sprintf("CPUutilization%d", i+1)] = x
			}
			mutex.Unlock()
		}
	}()

	go func() {
		for {
			<-pullSignalBasic

			if errorHappened {
				return
			}

			metrics := agent.ObtainMetricsSample()

			mutex.Lock()
			sample = metrics
			counter += 1
			mutex.Unlock()
		}
	}()

	go func() {
		for {
			pullSignalExtra <- struct{}{}
			pullSignalBasic <- struct{}{}
			time.Sleep(time.Duration(flags.Metrics.PollInterval) * time.Second)
		}
	}()

	s := semaphore.NewWeighted(int64(flags.Metrics.SendRateLimit))

	for {
		time.Sleep(time.Duration(flags.Metrics.RepeatInterval) * time.Second)

		s.Acquire(context.Background(), 1)
		go func() {
			mutex.Lock()
			metrics := sample.Clone().AddGauges(extraGauges).AddPollCount(counter)
			counter = 0
			mutex.Unlock()

			if len(metrics.Counters) == 0 && len(metrics.Gauges) == 0 {
				s.Release(1)
				return
			}

			resp, err := sendMetricsBatch(URLTemplate, flags.HashKey, metrics.Counters, metrics.Gauges)
			if err != nil || resp.StatusCode != 200 {
				if err != nil {
					log.Println(err.Error())
				} else {
					resp.Body.Close()
				}

				mutex.Lock()
				counter += metrics.PollCount()
				mutex.Unlock()
			} else {
				validateServerResponse(resp, flags.HashKey)

				resp.Body.Close()
			}
			s.Release(1)
		}()
	}
}

// No action is required per tech task, so the agent only logs the result if the response is not valid
func validateServerResponse(resp *http.Response, hashKey string) {
	if hashKey == "" {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to verify server response", err)
		return
	}
	correctHash := utils.Hash(body, hashKey)
	receivedHash := resp.Header.Get(utils.HashHeader)

	if correctHash != receivedHash {
		log.Println("Server response did not pass validation")
	}
}

func sendMetricsBatch(url string, hashKey string, counters map[string]int64, gauges map[string]float64) (*http.Response, error) {
	m := make([]models.Metrics, 0)
	for name, value := range counters {
		m = append(m, models.Metrics{
			ID:    name,
			MType: models.Counter,
			Delta: &value,
		})
	}
	for name, value := range gauges {
		m = append(m, models.Metrics{
			ID:    name,
			MType: models.Gauge,
			Value: &value,
		})
	}

	body, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var compressed bytes.Buffer
	gw, err := gzip.NewWriterLevel(&compressed, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	_, err = gw.Write(body)
	if err != nil {
		return nil, err
	}

	gw.Close()
	req, err := http.NewRequest(http.MethodPost, url, &compressed)
	if err != nil {
		return nil, err
	}
	err = utils.SetHashHeader(req, hashKey)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := utils.NewRetryingClient().Do(req)
	return resp, err
}

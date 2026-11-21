package main

import (
	"bytes"
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

	var pullSignalExtra = make(chan struct{}, 1)
	var pullSignalBasic = make(chan struct{}, 1)

	var childCtx = context.Background()
	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-childCtx.Done():
				return
			case <-pullSignalExtra:
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
		}
	}()

	go func() {
		defer wg.Done()

		for {
			select {
			case <-childCtx.Done():
				return
			case <-pullSignalBasic:
				metrics := agent.ObtainMetricsSample()

				mutex.Lock()
				sample = metrics
				counter += 1
				mutex.Unlock()
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-childCtx.Done():
				return
			case pullSignalExtra <- struct{}{}:
			}

			select {
			case <-childCtx.Done():
				return
			case pullSignalBasic <- struct{}{}:
			}
			time.Sleep(time.Duration(flags.Metrics.PollInterval) * time.Second)
		}
	}()

	jobs := make(chan agent.MetricsSample, flags.Metrics.SendRateLimit)
	results := make(chan sendWorkResult, flags.Metrics.SendRateLimit)

	// spawn workers
	var workersGroup sync.WaitGroup
	for w := 1; w <= flags.Metrics.SendRateLimit; w++ {
		workersGroup.Add(1)
		go func() {
			sendWorker(childCtx, URLTemplate, flags.HashKey, jobs, results)
			workersGroup.Done()
		}()
	}

	// wait for all workers to finish (in case of cancellation) and close the results channel
	go func() {
		defer wg.Done()

		workersGroup.Wait()
		close(results)
	}()

	// send work
	go func() {
		defer close(jobs)
		defer wg.Done()

		for {
			time.Sleep(time.Duration(flags.Metrics.RepeatInterval) * time.Second)

			mutex.Lock()
			metrics := sample.Clone().AddGauges(extraGauges).AddPollCount(counter)
			counter = 0
			mutex.Unlock()

			if len(metrics.Counters) == 0 && len(metrics.Gauges) == 0 {
				return
			}

			select {
			case <-childCtx.Done():
				return
			case jobs <- metrics:
			}
		}
	}()

	// recieve results and handle sending failure
	for result := range results {
		if !result.success {
			mutex.Lock()
			counter += result.pollCount
			mutex.Unlock()
		}
	}

	wg.Wait()
}

type sendWorkResult struct {
	success   bool
	pollCount int64
}

func sendWorker(ctx context.Context, URLTemplate string, hashKey string, jobs <-chan agent.MetricsSample, results chan<- sendWorkResult) {
	for metrics := range jobs {
		resp, err := sendMetricsBatch(ctx, URLTemplate, hashKey, metrics.Counters, metrics.Gauges)

		if err != nil || resp.StatusCode != 200 {
			if err != nil {
				log.Println(err.Error())
			} else {
				resp.Body.Close()
			}

			select {
			case <-ctx.Done():
				return
			case results <- sendWorkResult{
				success:   false,
				pollCount: metrics.PollCount(),
			}:
			}
		} else {
			validateServerResponse(resp, hashKey)

			resp.Body.Close()

			select {
			case <-ctx.Done():
				return
			case results <- sendWorkResult{
				success: true,
			}:
			}
		}
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

func sendMetricsBatch(ctx context.Context, url string, hashKey string, counters map[string]int64, gauges map[string]float64) (*http.Response, error) {
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := utils.NewSigningClient(utils.NewGzipClient(utils.NewRetryingClient(http.DefaultClient)), hashKey).Do(req)
	return resp, err
}

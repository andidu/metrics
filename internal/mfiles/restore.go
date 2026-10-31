package mfiles

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	"github.com/andidu/metrics/internal/handler"
	models "github.com/andidu/metrics/internal/model"
	"github.com/andidu/metrics/internal/service"
)

var ErrUnsupportedType = errors.New("unsupported metric type")
var ErrWrongFormat = errors.New("the metric doesn't have all required fields")

func Restore(filename string) (handler.MemStorage, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	memStorage := service.NewMemStorage()

	for {
		if !scanner.Scan() {
			err = scanner.Err()
			if err != nil {
				return nil, err
			} else {
				return memStorage, nil
			}
		}
		data := scanner.Bytes()

		var metric models.Metrics
		err = json.Unmarshal(data, &metric)
		if err != nil {
			return nil, err
		}
		if metric.MType == models.Gauge {
			if metric.Value != nil {
				memStorage.UpdateGauge(metric.ID, *metric.Value)
			} else {
				return nil, ErrWrongFormat
			}
		} else if metric.MType == models.Counter {
			if metric.Delta != nil {
				memStorage.OverrideCounter(metric.ID, *metric.Delta)
			} else {
				return nil, ErrWrongFormat
			}
		} else {
			return nil, ErrUnsupportedType
		}
	}
}

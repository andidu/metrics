package handler

import (
	"errors"
	"log"
	"strconv"
	"strings"

	models "github.com/andidu/metrics/internal/model"
)

type Service struct {
	storage MemStorage
}

var errNoElement = errors.New("no such element found")
var ErrUnknownMetricType = errors.New("unknown metric type")
var ErrWrongMetricValue = errors.New("wrong metric value")
var ErrInternalStorage = errors.New("internal storage error")

func (s Service) Gauges() map[string]float64 {
	return s.storage.Gauges()
}

func (s Service) Counters() map[string]int64 {
	return s.storage.Counters()
}

func (s Service) UpdateMetric(t string, name string, strvalue string) error {
	if t != "gauge" && t != "counter" {
		return ErrUnknownMetricType
	}
	if t == "gauge" {
		return s.updateGauge(name, strvalue)
	} else {
		return s.updateCounter(name, strvalue)
	}
}

func (s Service) updateGauge(name string, strvalue string) error {
	value, err := strconv.ParseFloat(strvalue, 64)
	if err != nil {
		return ErrWrongMetricValue
	}

	err = s.storage.UpdateGauge(name, value)
	if err != nil {
		log.Println("Internal storage error", err)
		return ErrInternalStorage
	}
	return nil
}

func (s Service) updateCounter(name string, strvalue string) error {
	value, err := strconv.Atoi(strvalue)
	if err != nil {
		return ErrWrongMetricValue
	}

	err = s.storage.UpdateCounter(name, value)
	if err != nil {
		log.Println("Internal storage error", err)
		return ErrInternalStorage
	}
	return nil
}

func (s Service) UpdateMetrics(m models.Metrics) error {
	t := m.MType
	if t != "gauge" && t != "counter" {
		return ErrUnknownMetricType
	}
	if t == "gauge" {
		err := s.storage.UpdateGauge(m.ID, *m.Value)
		if err != nil {
			log.Println("Internal storage error", err)
			return ErrInternalStorage
		}
		return nil
	} else {
		err := s.storage.UpdateCounter(m.ID, int(*m.Delta))
		if err != nil {
			log.Println("Internal storage error", err)
			return ErrInternalStorage
		}
		return nil
	}
}

func (s Service) GetGauge(name string) (string, error) {
	fval, ok := s.storage.Gauges()[name]
	if ok {
		s := strconv.FormatFloat(fval, 'f', 3, 64)
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
		return s, nil
	}
	return "", errNoElement
}
func (s Service) GetCounter(name string) (string, error) {
	ival, ok := s.storage.Counters()[name]
	if ok {
		return strconv.FormatInt(ival, 10), nil
	}
	return "", errNoElement
}

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

var ErrNoElement = errors.New("no such element found")
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
	value, err := strconv.ParseInt(strvalue, 10, 64)
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
		if m.Value == nil {
			return ErrWrongMetricValue
		}
		err := s.storage.UpdateGauge(m.ID, *m.Value)
		if err != nil {
			log.Println("Internal storage error", err)
			return ErrInternalStorage
		}
		return nil
	} else {
		if m.Delta == nil {
			return ErrWrongMetricValue
		}
		err := s.storage.UpdateCounter(m.ID, *m.Delta)
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
	return "", ErrNoElement
}
func (s Service) GetCounter(name string) (string, error) {
	ival, ok := s.storage.Counters()[name]
	if ok {
		return strconv.FormatInt(ival, 10), nil
	}
	return "", ErrNoElement
}

// Retreives metrics by type and name and sets its value into m
func (s Service) RetreiveMetrics(m *models.Metrics) error {
	t := m.MType
	if t != "gauge" && t != "counter" {
		return ErrUnknownMetricType
	}

	if t == "gauge" {
		fval, ok := s.storage.Gauges()[m.ID]
		if ok {
			m.Value = &fval
			return nil
		} else {
			return ErrNoElement
		}
	} else {
		ival, ok := s.storage.Counters()[m.ID]
		if ok {
			m.Delta = &ival
			return nil
		} else {
			return ErrNoElement
		}
	}
}

package handler

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

type Service struct {
	storage MemStorage
}

var noElementErr = errors.New("No such element found")
var UnknownMetricTypeErr = errors.New("Unknown metric type")
var WrongMetricValueErr = errors.New("Wrong metric value")
var InternalStorageErr = errors.New("Internal storage error")

func (s Service) Gauges() map[string]float64 {
	return s.storage.Gauges()
}

func (s Service) Counters() map[string]int64 {
	return s.storage.Counters()
}

func (s Service) UpdateMetric(t string, name string, strvalue string) error {
	if t != "gauge" && t != "counter" {
		return UnknownMetricTypeErr
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
		return WrongMetricValueErr
	}

	err = s.storage.UpdateGauge(name, value)
	if err != nil {
		log.Println("Internal storage error", err)
		return InternalStorageErr
	}
	return nil
}

func (s Service) updateCounter(name string, strvalue string) error {
	value, err := strconv.Atoi(strvalue)
	if err != nil {
		return WrongMetricValueErr
	}

	err = s.storage.UpdateCounter(name, value)
	if err != nil {
		log.Println("Internal storage error", err)
		return InternalStorageErr
	}
	return nil
}

func (s Service) GetGauge(name string) (string, error) {
	fval, ok := s.storage.Gauges()[name]
	if ok {
		s := strconv.FormatFloat(fval, 'f', 3, 64)
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
		return s, nil
	}
	return "", noElementErr
}
func (s Service) GetCounter(name string) (string, error) {
	ival, ok := s.storage.Counters()[name]
	if ok {
		return strconv.FormatInt(ival, 10), nil
	}
	return "", noElementErr
}

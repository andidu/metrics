package mfiles

import (
	"github.com/andidu/metrics/internal/handler"
)

type SavingStorage struct {
	handler.MemStorage
	Filename string
}

func (s SavingStorage) UpdateGauge(name string, value float64) error {
	err := s.MemStorage.UpdateGauge(name, value)
	if err == nil {
		return Save(s.Filename, s.MemStorage)
	}
	return err
}

func (s SavingStorage) UpdateCounter(name string, value int) error {
	err := s.MemStorage.UpdateCounter(name, value)
	if err == nil {
		return Save(s.Filename, s.MemStorage)
	}
	return err
}

func (s SavingStorage) OverrideCounter(name string, value int64) error {
	err := s.MemStorage.OverrideCounter(name, value)
	if err == nil {
		return Save(s.Filename, s.MemStorage)
	}
	return err
}

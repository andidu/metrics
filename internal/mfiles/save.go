package mfiles

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/andidu/metrics/internal/handler"
	models "github.com/andidu/metrics/internal/model"
)

func Save(filename string, m handler.MemStorage) error {
	// override file on every pass
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for name, value := range m.Counters() {
		m := models.Metrics{
			ID:    name,
			MType: models.Counter,
			Delta: &value,
		}
		data, err := json.Marshal(&m)
		if err != nil {
			return err
		}
		data = append(data, '\n')

		_, err = writer.Write(data)
		if err != nil {
			return err
		}
	}
	for name, value := range m.Gauges() {
		m := models.Metrics{
			ID:    name,
			MType: models.Gauge,
			Value: &value,
		}
		data, err := json.Marshal(&m)
		if err != nil {
			return err
		}
		data = append(data, '\n')

		_, err = writer.Write(data)
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

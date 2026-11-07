package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andidu/metrics/internal/handler"
	models "github.com/andidu/metrics/internal/model"
	"github.com/andidu/metrics/internal/router"
	"github.com/andidu/metrics/internal/service"
	"github.com/andidu/metrics/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_HandleUpdates(t *testing.T) {
	memStorage := service.NewMemStorage(nil)
	ts := httptest.NewServer(router.MetricsRouter(handler.New(memStorage)))
	defer ts.Close()

	var test2IValue = int64(9)
	var test1FValue = 9.01
	var test3IValue = int64(45)
	var test3FValue = 7.67

	metrics := []models.Metrics{
		models.Metrics{
			ID:    "test1",
			MType: models.Gauge,
			Value: &test1FValue,
			Delta: nil,
		},
		models.Metrics{
			ID:    "test2",
			MType: models.Counter,
			Value: nil,
			Delta: &test2IValue,
		},
		models.Metrics{
			ID:    "test3",
			MType: models.Gauge,
			Value: &test3FValue,
			Delta: nil,
		},
		models.Metrics{
			ID:    "test3",
			MType: models.Counter,
			Value: nil,
			Delta: &test3IValue,
		},
	}
	t.Run("Metrics Batch", func(t *testing.T) {
		body, err := json.Marshal(metrics)
		require.NoError(t, err)

		resp, _ := testutils.TestRequestWithBody(t, ts, http.MethodPost, "/updates/", bytes.NewReader(body))
		assert.Equal(t, 200, resp.StatusCode)

		gauges, err := memStorage.Gauges()
		assert.NoError(t, err)

		assert.Equal(t, test1FValue, gauges["test1"])
		assert.Equal(t, test3FValue, gauges["test3"])
		assert.NotContains(t, gauges, "test2")

		counters, err := memStorage.Counters()
		assert.NoError(t, err)

		assert.Equal(t, test2IValue, counters["test2"])
		assert.Equal(t, test3IValue, counters["test3"])
		assert.NotContains(t, counters, "test1")
		resp.Body.Close()
	})

}

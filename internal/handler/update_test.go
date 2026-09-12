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

func TestHandler_HandleUpdate(t *testing.T) {
	memStorage := service.NewMemStorage()
	ts := httptest.NewServer(router.MetricsRouter(handler.New(memStorage)))
	defer ts.Close()

	var v9 = int64(9)
	var v901 = 9.01

	tests := []struct {
		m models.Metrics
	}{
		{
			m: models.Metrics{
				ID:    "test1",
				MType: models.Gauge,
				Value: &v901,
				Delta: nil,
			},
		},
		{
			m: models.Metrics{
				ID:    "test2",
				MType: models.Counter,
				Value: nil,
				Delta: &v9,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.m.ID, func(t *testing.T) {
			body, err := json.Marshal(tt.m)
			require.NoError(t, err)

			var initialValue = int64(0)
			if tt.m.MType == models.Counter {
				initialValue = memStorage.Counters()[tt.m.ID]
			}

			resp, _ := testutils.TestRequestWithBody(t, ts, http.MethodPost, "/update/", bytes.NewReader(body))
			assert.Equal(t, 200, resp.StatusCode)
			if tt.m.MType == models.Gauge {
				assert.Equal(t, *tt.m.Value, memStorage.Gauges()[tt.m.ID])
			} else {
				assert.Equal(t, (*(tt.m.Delta))+initialValue, memStorage.Counters()[tt.m.ID])
			}
			resp.Body.Close()
		})
	}
}

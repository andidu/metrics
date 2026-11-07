package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andidu/metrics/internal/handler"
	"github.com/andidu/metrics/internal/router"
	"github.com/andidu/metrics/internal/service"
	"github.com/stretchr/testify/require"
)

type stubStorage struct {
	handler.MemStorage
	ping bool
}

func (s stubStorage) Ping() bool {
	return s.ping
}

func TestHandler_HandlePing(t *testing.T) {
	tests := []struct {
		ping         bool
		responseCode int
	}{
		{
			ping:         true,
			responseCode: 200,
		},
		{
			ping:         false,
			responseCode: 500,
		},
	}

	for _, tt := range tests {
		stub := stubStorage{
			MemStorage: service.NewMemStorage(nil),
			ping:       tt.ping,
		}

		ts := httptest.NewServer(router.MetricsRouter(handler.New(stub)))
		defer ts.Close()

		t.Run(fmt.Sprintf("ping %t", tt.ping), func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, ts.URL+"/ping", nil)
			require.NoError(t, err)

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tt.responseCode, resp.StatusCode)
		})
	}
}

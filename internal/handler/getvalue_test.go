package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andidu/metrics/internal/handler"
	"github.com/andidu/metrics/internal/service"
	"github.com/andidu/metrics/internal/testutils"
	"github.com/stretchr/testify/assert"

	"github.com/andidu/metrics/internal/router"
)

func TestHandleGetValueGauge(t *testing.T) {
	storage := service.NewMemStorage()
	ts := httptest.NewServer(router.MetricsRouter(handler.New(storage)))
	defer ts.Close()
	storage.Counters()["name"] = 5
	storage.Gauges()["name1"] = 9.011

	type want struct {
		contentType string
		statusCode  int
		body        string
	}
	type request struct {
		url    string
		method string
	}
	tests := []struct {
		name    string // description of this test case
		request request
		want    want
	}{
		{
			name: "happy path",
			request: request{
				url:    "/value/gauge/name1",
				method: http.MethodGet,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusOK,
				body:        "9.011",
			},
		},
		{
			name: "wrong method",
			request: request{
				url:    "/value/gauge/name1",
				method: http.MethodPost,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusMethodNotAllowed,
				body:        "",
			},
		},
		{
			name: "with value",
			request: request{
				url:    "/value/gauge/name1/6.901",
				method: http.MethodGet,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
				body:        "404 page not found\n",
			},
		},
		{
			name: "wrong url",
			request: request{
				url:    "/value/gauge/name1/name1",
				method: http.MethodGet,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
				body:        "404 page not found\n",
			},
		},
		{
			name: "wrong name",
			request: request{
				url:    "/value/gauge/name",
				method: http.MethodGet,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
				body:        "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, get := testutils.TestRequest(t, ts, tt.request.method, tt.request.url)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.body, get)

		})
	}
}

func TestHandleGetValueCounter(t *testing.T) {
	storage := service.NewMemStorage()
	ts := httptest.NewServer(router.MetricsRouter(handler.New(storage)))
	defer ts.Close()
	storage.Counters()["name"] = 5
	storage.Gauges()["name1"] = 9.011

	type want struct {
		contentType string
		statusCode  int
		body        string
	}
	type request struct {
		url    string
		method string
	}
	tests := []struct {
		name    string // description of this test case
		request request
		want    want
	}{
		{
			name: "happy path",
			request: request{
				url:    "/value/counter/name",
				method: http.MethodGet,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusOK,
				body:        "5",
			},
		},
		{
			name: "wrong method",
			request: request{
				url:    "/value/counter/name",
				method: http.MethodPost,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusMethodNotAllowed,
				body:        "",
			},
		},
		{
			name: "with value",
			request: request{
				url:    "/value/counter/name/9",
				method: http.MethodGet,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
				body:        "404 page not found\n",
			},
		},
		{
			name: "wrong url",
			request: request{
				url:    "/value/counter/name/name",
				method: http.MethodGet,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
				body:        "404 page not found\n",
			},
		},
		{
			name: "wrong name",
			request: request{
				url:    "/value/counter/name1",
				method: http.MethodGet,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
				body:        "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, get := testutils.TestRequest(t, ts, tt.request.method, tt.request.url)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.body, get)

		})
	}
}

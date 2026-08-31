package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andidu/metrics/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestHandleUpdateGauge(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
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
				url:    "/update/gauge/name/4.001",
				method: http.MethodPost,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
		},
		{
			name: "wrong method",
			request: request{
				url:    "/update/gauge/name/4.001",
				method: http.MethodGet,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusMethodNotAllowed,
			},
		},
		{
			name: "no value",
			request: request{
				url:    "/update/gauge/name",
				method: http.MethodPost,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name: "wrong url",
			request: request{
				url:    "/update/gauge/name/name/4.001",
				method: http.MethodPost,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name: "wrong value",
			request: request{
				url:    "/update/gauge/name/4.o01",
				method: http.MethodPost,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.url, nil)

			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(handler.HandleUpdate)
			handler(recorder, request)

			result := recorder.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}

func TestHandleUpdateCounter(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
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
				url:    "/update/counter/name/4",
				method: http.MethodPost,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
		},
		{
			name: "wrong method",
			request: request{
				url:    "/update/counter/name/4",
				method: http.MethodGet,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusMethodNotAllowed,
			},
		},
		{
			name: "no value",
			request: request{
				url:    "/update/counter/name",
				method: http.MethodPost,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name: "wrong url",
			request: request{
				url:    "/update/counter/name/name/4",
				method: http.MethodPost,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name: "wrong value",
			request: request{
				url:    "/update/counter/name/4.o01",
				method: http.MethodPost,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.url, nil)

			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(handler.HandleUpdate)
			handler(recorder, request)

			result := recorder.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}

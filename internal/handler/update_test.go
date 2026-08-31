package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andidu/metrics/internal/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(
	t *testing.T,
	ts *httptest.Server,
	method, path string,
) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestHandleUpdateGauge(t *testing.T) {
	ts := httptest.NewServer(router.MetricsRouter())
	defer ts.Close()

	type want struct {
		contentType string
		statusCode  int
		emptyBody   bool
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
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusOK,
				emptyBody:   true,
			},
		},
		{
			name: "wrong method",
			request: request{
				url:    "/update/gauge/name/4.001",
				method: http.MethodGet,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusMethodNotAllowed,
				emptyBody:   true,
			},
		},
		{
			name: "no value",
			request: request{
				url:    "/update/gauge/name",
				method: http.MethodPost,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
				emptyBody:   false,
			},
		},
		{
			name: "wrong url",
			request: request{
				url:    "/update/gauge/name/name/4.001",
				method: http.MethodPost,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
				emptyBody:   false,
			},
		},
		{
			name: "wrong value",
			request: request{
				url:    "/update/gauge/name/4.o01",
				method: http.MethodPost,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				emptyBody:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, get := testRequest(t, ts, tt.request.method, tt.request.url)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			if tt.want.emptyBody {
				assert.Equal(t, "", get)
			}
		})
	}
}

func TestHandleUpdateCounter(t *testing.T) {
	ts := httptest.NewServer(router.MetricsRouter())
	defer ts.Close()

	type want struct {
		contentType string
		statusCode  int
		emptyBody   bool
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
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusOK,
				emptyBody:   true,
			},
		},
		{
			name: "wrong method",
			request: request{
				url:    "/update/counter/name/4",
				method: http.MethodGet,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusMethodNotAllowed,
				emptyBody:   true,
			},
		},
		{
			name: "no value",
			request: request{
				url:    "/update/counter/name",
				method: http.MethodPost,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
				emptyBody:   true,
			},
		},
		{
			name: "wrong url",
			request: request{
				url:    "/update/counter/name/name/4",
				method: http.MethodPost,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
				emptyBody:   false,
			},
		},
		{
			name: "wrong value",
			request: request{
				url:    "/update/counter/name/4.o01",
				method: http.MethodPost,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				emptyBody:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, get := testRequest(t, ts, tt.request.method, tt.request.url)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			if tt.want.emptyBody {
				assert.Equal(t, "", get)
			}
		})
	}
}

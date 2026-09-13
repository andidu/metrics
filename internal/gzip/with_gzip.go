package gzip

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/andidu/metrics/internal/utils"
)

func WithGzip(h http.Handler) http.Handler {
	gzipFunc := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		gzw, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			utils.Logger.Errorln("Unable to create new gzip writer")
			return
		}

		defer gzw.Close()

		w.Header().Set("Content-Encoding", "gzip")
		h.ServeHTTP(GzipWriter{
			ResponseWriter: w,
			W:              gzw,
		}, r)
	}

	return http.HandlerFunc(gzipFunc)
}

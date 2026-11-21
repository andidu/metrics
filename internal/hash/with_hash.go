package hash

import (
	"net/http"

	"github.com/andidu/metrics/internal/utils"
)

func WithHash(handler http.Handler, hashKey string) http.Handler {
	fun := func(w http.ResponseWriter, r *http.Request) {
		ok, err := utils.CheckHashHeader(r, hashKey)

		if err != nil {
			utils.Logger.Errorln("Unable to check hash header:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if hashKey == "" {
			handler.ServeHTTP(w, r)
			return
		}

		hashWriter := HashWriter{
			ResponseWriter: w,
			BuffHeader:     http.Header{},
			BuffStatusCode: 200,
			BuffResponse:   []byte{},
		}

		handler.ServeHTTP(&hashWriter, r)

		hash := utils.Hash(hashWriter.BuffResponse, hashKey)
		hashWriter.BuffHeader.Set(utils.HashHeader, hash)

		for name, values := range hashWriter.BuffHeader {
			w.Header()[name] = values
		}
		hashWriter.ResponseWriter.WriteHeader(hashWriter.BuffStatusCode)
		hashWriter.ResponseWriter.Write(hashWriter.BuffResponse)
	}

	return http.HandlerFunc(fun)
}

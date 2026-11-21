package utils

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

type GzipClient struct {
	Client HTTPClient
}

func NewGzipClient(client HTTPClient) GzipClient {
	return GzipClient{
		Client: client,
	}
}

func (uc GzipClient) Do(req *http.Request) (*http.Response, error) {
	compressedReq, err := compressRequest(req)
	if err != nil {
		return nil, err
	}

	resp, err := uc.Client.Do(compressedReq)
	if err != nil {
		return nil, err
	}

	return decompressResponse(resp)
}

func decompressResponse(resp *http.Response) (*http.Response, error) {
	if resp.Header.Get("Content-Encoding") == "gzip" {
		defer resp.Body.Close()

		gzr, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer gzr.Close()

		body, err := io.ReadAll(gzr)
		if err != nil {
			return nil, err
		}

		// Clean up the response:
		resp.Header.Del("Content-Encoding")
		resp.Header.Del("Content-Length")
		resp.ContentLength = int64(len(body))
		resp.Uncompressed = true

		resp.Body = io.NopCloser(bytes.NewReader(body))
	}

	return resp, nil
}

func compressRequest(req *http.Request) (*http.Request, error) {
	var compressed bytes.Buffer
	gw, err := gzip.NewWriterLevel(&compressed, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	_, err = gw.Write(body)
	if err != nil {
		return nil, err
	}

	err = gw.Close()
	if err != nil {
		return nil, err
	}

	compressedReq, err := http.NewRequestWithContext(
		req.Context(),
		req.Method,
		req.URL.String(),
		bytes.NewReader(compressed.Bytes()),
	)
	if err != nil {
		return nil, err
	}

	compressedReq.Header = req.Header.Clone()
	compressedReq.Header.Set("Content-Encoding", "gzip")

	return compressedReq, nil
}

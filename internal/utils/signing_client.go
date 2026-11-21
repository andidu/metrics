package utils

import (
	"net/http"
)

type SigningClient struct {
	Client  HttpClient
	hashKey string
}

func NewSigningClient(hc HttpClient, hashKey string) HttpClient {
	return SigningClient{
		Client:  hc,
		hashKey: hashKey,
	}
}

func (c SigningClient) Do(req *http.Request) (*http.Response, error) {
	err := SetHashHeader(req, c.hashKey)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

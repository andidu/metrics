package utils

import (
	"errors"
	"net/http"
	"syscall"
)

type RetryingClient struct {
	*http.Client
}

func NewRetryingClient() RetryingClient {
	return RetryingClient{
		http.DefaultClient,
	}
}

func (c RetryingClient) Do(req *http.Request) (*http.Response, error) {
	firstAttempt := true
	hasBody := req.Body != nil && req.Body != http.NoBody

	return TryValue(
		func() (*http.Response, error) {
			if firstAttempt {
				firstAttempt = false
				return c.Client.Do(req)
			}

			newReq := req.Clone(req.Context())

			if hasBody {
				if newReq.GetBody == nil {
					return nil, errors.New("can't retry request: body can't be recreated")
				}

				body, err := newReq.GetBody()

				if err != nil {
					return nil, err
				}

				newReq.Body = body
			}

			return c.Client.Do(newReq)

		},
		func(err error) bool {
			return errors.Is(err, syscall.ECONNREFUSED)
		},
	)
}

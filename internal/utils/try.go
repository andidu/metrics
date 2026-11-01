package utils

import "time"

func Try(operation func() error, isRetryable func(error) bool) error {
	_, err := TryValue(
		func() (struct{}, error) {
			return struct{}{}, operation()
		},
		isRetryable,
	)
	return err
}

func TryValue[T any](operation func() (T, error), isRetryable func(error) bool) (T, error) {
	delays := []time.Duration{
		time.Second,
		3 * time.Second,
		5 * time.Second,
	}

	var err error
	var t T

	for attempt := 0; attempt <= len(delays); attempt++ {
		t, err = operation()

		if err == nil {
			return t, nil
		}
		if !isRetryable(err) {
			return t, err
		}

		if attempt < len(delays) {
			time.Sleep(delays[attempt])
		}
	}

	return t, err
}

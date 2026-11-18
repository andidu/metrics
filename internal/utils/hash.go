package utils

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"slices"
)

const HashHeader = "HashSHA256"

func SetHashHeader(req *http.Request, hashKey string) error {
	if hashKey == "" {
		return nil
	}

	content, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(content))

	bodyhHash := Hash(content, hashKey)
	req.Header.Set(HashHeader, bodyhHash)
	return nil
}

// Returns `true` if hash value is valid or not present, `false` otherwise
func CheckHashHeader(req *http.Request, hashKey string) (bool, error) {
	if hashKey == "" {
		return true, nil
	}

	receivedHash := req.Header.Get(HashHeader)

	content, err := io.ReadAll(req.Body)
	if err != nil {
		return false, err
	}

	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(content))

	bodyhHash := Hash(content, hashKey)
	return bodyhHash == receivedHash, nil
}

func Hash(content []byte, hashKey string) string {
	content = slices.Concat(content, []byte(hashKey))

	sum := sha256.Sum256(content)
	return fmt.Sprintf("%x", sum)
}

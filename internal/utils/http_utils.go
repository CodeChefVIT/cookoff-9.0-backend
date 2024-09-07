package utils

import (
	"io"
	"net/http"
)

// DoRequest handles making HTTP requests and returning the response body.
func DoRequest(method, url string, body io.Reader, client *http.Client) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

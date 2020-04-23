package tool

import (
	"fmt"
	"io"
	"net/http"
)

func Get(url string, hClient http.Client) (io.ReadCloser, error) {
	resp, err := hClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http error: status code %d", resp.StatusCode)
	}
	return resp.Body, nil
}

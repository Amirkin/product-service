package main

import (
	"errors"
	"io"
	"net/http"
)

type Requester struct {
	client http.Client
}

func (r *Requester) GetCSV(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return nil, errors.New("body is nil")
	}

	return resp.Body, nil
}

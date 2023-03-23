package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type HttpResponse struct {
	StatusCode int
	Message    string
	StatusText string
}

type Sender struct {
	client    *http.Client
	serverUrl string
	apiKey    string
}

func NewHttpClient() *Sender {
	return &Sender{
		serverUrl: "https://atcqt90a7j.execute-api.ap-southeast-1.amazonaws.com/POC-DDB-TiDB-DevStage/shipments",
		apiKey:    "H9Ze4AzmXz1yXW1uVZYs31soaD3JHOT95gjp8QjT",
		client:    &http.Client{},
	}
}

func (s *Sender) Post(data []byte) (HttpResponse, error) {
	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, s.serverUrl, bytes.NewBuffer(data))
	if err != nil {
		return HttpResponse{StatusCode: 500, StatusText: "Internal Error", Message: "failed to create new request"}, err
	}

	request.Header.Set("x-api-key", s.apiKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := s.client.Do(request)
	if err != nil {
		return HttpResponse{StatusCode: response.StatusCode, StatusText: response.Status, Message: "failed to do request"}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return HttpResponse{StatusCode: response.StatusCode, StatusText: response.Status, Message: string(body)}, err
	}

	return HttpResponse{StatusCode: response.StatusCode, StatusText: response.Status, Message: string(body)}, err
}

package http

import "net/http"

type Service struct {
	movieClient  *http.Client
	movieBaseURL string
}

func New(client *http.Client, baseURL string) *Service {
	return &Service{
		movieClient:  client,
		movieBaseURL: baseURL,
	}
}

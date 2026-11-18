package adapter

import "net/http"

type MovieClient struct {
	client       *http.Client
	movieBaseURL string
}

func New(client *http.Client, movieBaseURL string) *MovieClient {
	return &MovieClient{
		client:       client,
		movieBaseURL: movieBaseURL,
	}
}

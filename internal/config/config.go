package config

import "os"

type Config struct {
	MovieBaseURL string
	HTTPAddr     string
}

func LoadConfig() Config {
	movieBaseURL := os.Getenv("MOVIE_BASE_URL")
	if movieBaseURL == "" {
		movieBaseURL = "https://homeworksite.site"
	}

	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080"
	}

	return Config{
		MovieBaseURL: movieBaseURL,
		HTTPAddr:     httpAddr,
	}
}

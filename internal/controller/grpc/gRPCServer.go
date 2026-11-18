package grpc

import (
	"net/http"

	"github.com/arinamklvch/go-courses/internal/pb"
)

type GRPCServer struct {
	pb.UnimplementedMoviesServiceServer
	movieClient  *http.Client
	movieBaseURL string
}

func New(client *http.Client, baseURL string) *GRPCServer {
	return &GRPCServer{
		movieClient:  client,
		movieBaseURL: baseURL,
	}
}

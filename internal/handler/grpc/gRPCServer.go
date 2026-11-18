package grpc

import (
	"github.com/arinamklvch/go-courses/internal/pb"
	"github.com/arinamklvch/go-courses/internal/usecase"
)

type GRPCServer struct {
	pb.UnimplementedMoviesServiceServer
	movieUseCase usecase.MovieUseCase
}

func New(movieUseCase usecase.MovieUseCase) *GRPCServer {
	return &GRPCServer{
		movieUseCase: movieUseCase,
	}
}

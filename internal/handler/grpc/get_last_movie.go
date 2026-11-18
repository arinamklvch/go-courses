package grpc

import (
	"context"
	"errors"
	"log"

	"github.com/arinamklvch/go-courses/internal/domain"
	"github.com/arinamklvch/go-courses/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) GetLastMovie(ctx context.Context, _ *pb.GetLastMovieRequest) (*pb.Movie, error) {
	lastMovie, err := s.movieUseCase.GetLastMovie(ctx)
	if ctx.Err() != nil {
		log.Printf("gRPC GetLastMovie: request ended: %v", ctx.Err())
		return nil, status.FromContextError(ctx.Err()).Err()
	}

	if err != nil {
		log.Printf("gRPC GetLastMovie: %v", err)
		if errors.Is(err, domain.ErrMovieNotFound) {
			return nil, status.Error(codes.NotFound, "no movies found")
		}
		code := codes.Unavailable
		if errors.Is(err, domain.ErrCreateMovieRequest) {
			code = codes.Internal
		}
		return nil, status.Error(code, "failed to find last movie")
	}

	return &pb.Movie{
		Id:       int64(lastMovie.ID),
		Title:    lastMovie.Title,
		Year:     int64(lastMovie.Year),
		Director: lastMovie.Director,
	}, nil
}

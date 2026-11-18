package grpc

import (
	"context"
	"log"

	"github.com/arinamklvch/go-courses/internal/movies"
	"github.com/arinamklvch/go-courses/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) LoadMovies(ctx context.Context, rangeRequest *pb.LoadMoviesRequest) (*pb.LoadMoviesResponse, error) {
	client := s.movieClient

	from := int(rangeRequest.GetFrom())
	to := int(rangeRequest.GetTo())

	if from <= 0 {
		return nil, status.Error(codes.InvalidArgument, "parameter 'from' must be a positive integer")
	}

	if from > to {
		return nil, status.Error(codes.InvalidArgument, "parameter 'from' must be less than or equal to 'to'")
	}

	loadMoviesRangeResult, successCount := movies.GetRange(ctx, from, to, s.movieBaseURL, client)
	if ctx.Err() != nil {
		log.Printf("gRPC LoadMovies (from=%d, to=%d): request ended: %v", from, to, ctx.Err())
		return nil, status.FromContextError(ctx.Err()).Err()
	} else {
		log.Printf("gRPC LoadMovies (from=%d, to=%d): loaded %d of %d movies", from, to, successCount, to-from+1)
	}

	if successCount == 0 {
		return nil, status.Error(codes.Unavailable, "failed to load any movies")
	}

	results := make([]*pb.MovieResult, 0, len(loadMoviesRangeResult))
	for _, result := range loadMoviesRangeResult {
		movieResult := &pb.MovieResult{
			Id:    int64(result.ID),
			Error: result.Error,
		}
		if result.Movie != nil {
			movieResult.Movie = &pb.Movie{
				Id:       int64(result.Movie.ID),
				Title:    result.Title,
				Year:     int64(result.Year),
				Director: result.Director,
			}
		}
		results = append(results, movieResult)
	}

	return &pb.LoadMoviesResponse{Results: results}, nil
}

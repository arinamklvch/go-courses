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

func (s *GRPCServer) LoadMovies(ctx context.Context, rangeRequest *pb.LoadMoviesRequest) (*pb.LoadMoviesResponse, error) {
	from := int(rangeRequest.GetFrom())
	to := int(rangeRequest.GetTo())

	loadMoviesRangeResult, successCount, err := s.movieUseCase.GetMoviesRange(ctx, from, to)
	if ctx.Err() != nil {
		log.Printf("gRPC LoadMovies (from=%d, to=%d): request ended: %v", from, to, ctx.Err())
		return nil, status.FromContextError(ctx.Err()).Err()
	}
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrIncorrectID):
			return nil, status.Error(codes.InvalidArgument, "parameter 'from' must be a positive integer")
		case errors.Is(err, domain.ErrInvalidRange):
			return nil, status.Error(codes.InvalidArgument, "parameter 'from' must be less than or equal to 'to'")
		default:
			log.Printf("gRPC LoadMovies: %v", err)
			return nil, status.Error(codes.Internal, "failed to load movies")
		}
	}
	log.Printf("gRPC LoadMovies (from=%d, to=%d): loaded %d of %d movies", from, to, successCount, to-from+1)

	// куда????????????????????????????????????????????????????????????????
	if successCount == 0 {
		return nil, status.Error(codes.Unavailable, "failed to load any movies")
	}

	loadMoviesRangeResponse := make([]*pb.MovieResult, 0, len(loadMoviesRangeResult))
	for _, movieResult := range loadMoviesRangeResult {
		if movieResult.Error != nil {
			loadMoviesRangeResponse = append(loadMoviesRangeResponse, &pb.MovieResult{
				Id:    int64(movieResult.ID),
				Error: movieResult.Error.Error(),
			})
			continue
		}
		pbMovie := pb.Movie{
			Id:       int64(movieResult.ID),
			Title:    movieResult.Movie.Title,
			Year:     int64(movieResult.Movie.Year),
			Director: movieResult.Movie.Director,
		}
		loadMoviesRangeResponse = append(loadMoviesRangeResponse, &pb.MovieResult{
			Id:    int64(movieResult.ID),
			Movie: &pbMovie,
		})
	}

	return &pb.LoadMoviesResponse{Results: loadMoviesRangeResponse}, nil
}

package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/arinamklvch/go-courses/internal/domain"
)

func (s *Service) GetLastMovie(ctx context.Context) (domain.Movie, error) {
	lastFilmID, err := FindLastMovieID(ctx, s.MovieProvider)
	if err != nil {
		return domain.Movie{}, err
	}

	lastMovie, err := s.MovieProvider.GetMovie(ctx, lastFilmID)
	if err != nil {
		return domain.Movie{}, err
	}

	return lastMovie, nil
}

func FindLastMovieID(ctx context.Context, movieProvider MovieProvider) (int, error) {
	leftID := 1
	rightID := 100000
	var lastFilmID int

	for {
		_, err := movieProvider.GetMovie(ctx, rightID)
		if err == nil {
			leftID = rightID
			rightID *= 2
			continue
		}
		if errors.Is(err, domain.ErrMovieNotFound) {
			break
		}
		return 0, err
	}

	for leftID <= rightID {
		midID := (leftID + rightID) / 2
		_, err := movieProvider.GetMovie(ctx, midID)
		if err != nil {
			if errors.Is(err, domain.ErrMovieNotFound) {
				rightID = midID - 1
				continue
			}
			return 0, fmt.Errorf("movie %d: check availability: %w", midID, err)
		}
		lastFilmID = midID
		leftID = midID + 1
	}

	if lastFilmID == 0 {
		return 0, domain.ErrMovieNotFound
	}

	return lastFilmID, nil
}

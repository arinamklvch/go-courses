package usecase

import (
	"context"

	"github.com/arinamklvch/go-courses/internal/domain"
)

type MovieUseCase interface {
	GetMoviesRange(context.Context, int, int) ([]domain.MovieResult, int, error)
	GetLastMovie(context.Context) (domain.Movie, error)
}

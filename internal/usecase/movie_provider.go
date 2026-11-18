package usecase

import (
	"context"

	"github.com/arinamklvch/go-courses/internal/domain"
)

type MovieProvider interface {
	GetMovie(ctx context.Context, id int) (domain.Movie, error)
}

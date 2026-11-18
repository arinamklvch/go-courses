package http

import (
	"github.com/arinamklvch/go-courses/internal/usecase"
)

type Handler struct {
	movieUseCase usecase.MovieUseCase
}

func New(movieUseCase usecase.MovieUseCase) *Handler {
	return &Handler{
		movieUseCase: movieUseCase,
	}
}

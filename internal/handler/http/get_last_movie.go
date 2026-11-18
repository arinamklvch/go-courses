package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/arinamklvch/go-courses/internal/domain"
	"github.com/arinamklvch/go-courses/internal/dto"
)

// @Summary Get the last available movie
// @Description Finds the highest available movie ID by expanding the search bounds and using binary search. Assumes available movie IDs are consecutive starting at 1.
// @Tags movies
// @Produce json,plain
// @Success 200 {object} dto.Movie "The last available movie"
// @Failure 404 {string} string "No movies found; plain-text error"
// @Failure 500 {string} string "Failed to create a request to the movie service; plain-text error"
// @Failure 502 {string} string "Failed to search for or download the movie; plain-text error"
// @Router /last-movie [get]
func (h *Handler) GetLastMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	lastMovieResult, err := h.movieUseCase.GetLastMovie(ctx)
	if ctx.Err() != nil {
		log.Printf("GET /last-movie: request ended: %v", ctx.Err())
		return
	}
	if err != nil {
		log.Printf("GET /last-movie: %v", err)
		if errors.Is(err, domain.ErrMovieNotFound) {
			http.Error(w, "no movies found", http.StatusNotFound)
			return
		}
		status := http.StatusBadGateway
		if errors.Is(err, domain.ErrCreateMovieRequest) {
			status = http.StatusInternalServerError
		}
		http.Error(w, "failed to find last movie", status)
		return
	}

	lastMovieResponse := dto.Movie{
		ID:       lastMovieResult.ID,
		Title:    lastMovieResult.Title,
		Year:     lastMovieResult.Year,
		Director: lastMovieResult.Director,
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(lastMovieResponse)
	if err != nil {
		log.Printf("GET /last-movie: write JSON response: %v", err)
	}
}

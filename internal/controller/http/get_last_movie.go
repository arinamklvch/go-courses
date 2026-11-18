package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/arinamklvch/go-courses/internal/movies"
)

// @Summary Get the last available movie
// @Description Searches movie IDs from 1 to 200000 and returns the movie with the highest available ID in that range.
// @Tags movies
// @Produce json,plain
// @Success 200 {object} data.Movie "The last available movie in the search range"
// @Failure 404 {string} string "No movies found in the search range; plain-text error"
// @Failure 500 {string} string "Failed to create a request to the movie service; plain-text error"
// @Failure 502 {string} string "Failed to search for or download the movie; plain-text error"
// @Router /last-movie [get]
func (s *Service) GetLastMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client := s.movieClient

	lastMovie, err := movies.GetLast(ctx, client, s.movieBaseURL)
	if ctx.Err() != nil {
		log.Printf("GET /last-movie: request ended: %v", ctx.Err())
		return
	}
	if err != nil {
		log.Printf("GET /last-movie: %v", err)
		if errors.Is(err, movies.ErrMovieNotFound) {
			http.Error(w, "no movies found", http.StatusNotFound)
			return
		}
		status := http.StatusBadGateway
		if errors.Is(err, movies.ErrCreateMovieRequest) {
			status = http.StatusInternalServerError
		}
		http.Error(w, "failed to find last movie", status)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(lastMovie)
	if err != nil {
		log.Printf("GET /last-movie: write JSON response: %v", err)
	}
}

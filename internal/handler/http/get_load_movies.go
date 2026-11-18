package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/arinamklvch/go-courses/internal/domain"
	"github.com/arinamklvch/go-courses/internal/dto"
)

// @Summary Load a range of movies
// @Description Loads the inclusive range [from, to]. Each result contains the requested ID and either a movie or an error. Partial failures are included in the 200 response; 502 means no movies were loaded. Results may arrive in any order.
// @Tags movies
// @Produce json,plain
// @Param from query int true "First movie ID (inclusive)" minimum(1)
// @Param to query int true "Last movie ID (inclusive); must be greater than or equal to from" minimum(1)
// @Success 200 {array} dto.MovieResponse "At least one movie loaded; individual errors may also be present"
// @Failure 400 {string} string "Missing or invalid parameters, or from is greater than to; plain-text error"
// @Failure 500 {string} string "Unexpected use case error; plain-text error"
// @Failure 502 {array} dto.MovieResponse "No movies loaded; the array contains an error for each requested ID"
// @Router /load-movies [get]
func (h *Handler) GetMoviesRange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	from, err := strconv.Atoi(r.URL.Query().Get("from"))
	if err != nil {
		http.Error(w, "query parameter 'from' must be a positive integer", http.StatusBadRequest)
		return
	}

	to, err := strconv.Atoi(r.URL.Query().Get("to"))
	if err != nil {
		http.Error(w, "query parameter 'to' must be an integer", http.StatusBadRequest)
		return
	}

	loadMoviesRangeResult, successCount, err := h.movieUseCase.GetMoviesRange(ctx, from, to)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrIncorrectID):
			http.Error(w, "query parameter 'from' must be a positive integer", http.StatusBadRequest)
		case errors.Is(err, domain.ErrInvalidRange):
			http.Error(w, "query parameter 'from' must be less than or equal to 'to'", http.StatusBadRequest)
		default:
			log.Printf("GET /load-movies: %v", err)
			http.Error(w, "failed to load movies", http.StatusInternalServerError)
		}
		return
	}

	if ctx.Err() != nil {
		log.Printf("GET /load-movies (from=%d, to=%d): request ended: %v", from, to, ctx.Err())
		return
	} else {
		log.Printf("GET /load-movies (from=%d, to=%d): loaded %d of %d movies", from, to, successCount, to-from+1)
	}

	loadMoviesRangeResponse := make([]dto.MovieResponse, 0, len(loadMoviesRangeResult))
	for _, movieResult := range loadMoviesRangeResult {
		if movieResult.Error != nil {
			loadMoviesRangeResponse = append(loadMoviesRangeResponse, dto.MovieResponse{
				ID:    movieResult.ID,
				Error: movieResult.Error.Error(),
			})
			continue
		}
		dtoMovie := dto.Movie{
			ID:       movieResult.ID,
			Title:    movieResult.Movie.Title,
			Year:     movieResult.Movie.Year,
			Director: movieResult.Movie.Director,
		}
		loadMoviesRangeResponse = append(loadMoviesRangeResponse, dto.MovieResponse{
			ID:    movieResult.ID,
			Movie: &dtoMovie,
		})
	}

	w.Header().Set("Content-Type", "application/json")

	if successCount == 0 {
		w.WriteHeader(http.StatusBadGateway)
	}

	err = json.NewEncoder(w).Encode(loadMoviesRangeResponse)
	if err != nil {
		log.Printf("GET /load-movies (from=%d, to=%d): write JSON response: %v", from, to, err)
	}
}

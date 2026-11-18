package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/arinamklvch/go-courses/internal/domain"
	"github.com/arinamklvch/go-courses/internal/dto"
)

func (c *MovieClient) GetMovie(ctx context.Context, filmID int) (domain.Movie, error) {
	url := fmt.Sprintf("%s/%d/info.0.json", strings.TrimRight(c.movieBaseURL, "/"), filmID)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.Movie{}, fmt.Errorf("movie %d: %w: %w", filmID, domain.ErrCreateMovieRequest, err)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return domain.Movie{}, fmt.Errorf("movie %d: send HTTP request: %w", filmID, err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusNotFound {
		return domain.Movie{}, fmt.Errorf("movie %d: %w", filmID, domain.ErrMovieNotFound)
	}
	if response.StatusCode != http.StatusOK {
		return domain.Movie{}, fmt.Errorf("movie %d: upstream returned HTTP %d", filmID, response.StatusCode)
	}

	var movie dto.Movie
	if err := json.NewDecoder(response.Body).Decode(&movie); err != nil {
		return domain.Movie{}, fmt.Errorf("movie %d: decode JSON response: %w", filmID, err)
	}

	return domain.Movie{
		ID:       movie.ID,
		Title:    movie.Title,
		Year:     movie.Year,
		Director: movie.Director,
	}, nil
}

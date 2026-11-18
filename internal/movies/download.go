package movies

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/arinamklvch/go-courses/internal/data"
)

func Download(ctx context.Context, client *http.Client, movieBaseURL string, filmID int) (data.Movie, error) {
	url := fmt.Sprintf("%s/%d/info.0.json", strings.TrimRight(movieBaseURL, "/"), filmID)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return data.Movie{}, fmt.Errorf("movie %d: %w: %w", filmID, ErrCreateMovieRequest, err)
	}

	response, err := client.Do(request)
	if err != nil {
		return data.Movie{}, fmt.Errorf("movie %d: send HTTP request: %w", filmID, err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusNotFound {
		return data.Movie{}, fmt.Errorf("movie %d: %w", filmID, ErrMovieNotFound)
	}
	if response.StatusCode != http.StatusOK {
		return data.Movie{}, fmt.Errorf("movie %d: upstream returned HTTP %d", filmID, response.StatusCode)
	}

	var movie data.Movie
	if err := json.NewDecoder(response.Body).Decode(&movie); err != nil {
		return data.Movie{}, fmt.Errorf("movie %d: decode JSON response: %w", filmID, err)
	}

	return movie, nil
}

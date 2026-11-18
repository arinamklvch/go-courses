package movies

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/arinamklvch/go-courses/internal/data"
)

func GetLast(ctx context.Context, client *http.Client, movieBaseURL string) (data.Movie, error) {
	lastFilmID, err := FindLastMovieID(ctx, client, movieBaseURL)
	if err != nil {
		return data.Movie{}, err
	}

	lastMovie, err := Download(ctx, client, movieBaseURL, lastFilmID)
	if err != nil {
		return data.Movie{}, err
	}

	return lastMovie, nil
}

func FindLastMovieID(ctx context.Context, client *http.Client, movieBaseURL string) (int, error) {
	leftID := 1
	rightID := 100000
	var lastFilmID int

	for {
		_, err := Download(ctx, client, movieBaseURL, rightID)
		if err == nil {
			rightID *= 2
			continue
		}
		if errors.Is(err, ErrMovieNotFound) {
			break
		}
		return 0, err
	}

	for leftID <= rightID {
		midID := (leftID + rightID) / 2
		url := fmt.Sprintf("%s/%d/info.0.json", strings.TrimRight(movieBaseURL, "/"), midID)
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return 0, fmt.Errorf("movie %d: %w: %w", midID, ErrCreateMovieRequest, err)
		}
		response, err := client.Do(request)
		if err != nil {
			return 0, fmt.Errorf("movie %d: check availability: %w", midID, err)
		}
		_ = response.Body.Close()
		if response.StatusCode != http.StatusOK {
			if response.StatusCode == http.StatusNotFound {
				rightID = midID - 1
				continue
			}
			return 0, fmt.Errorf("movie %d: movie service returned HTTP %d", midID, response.StatusCode)
		} else {
			lastFilmID = midID
			leftID = midID + 1
		}
	}

	if lastFilmID == 0 {
		return 0, ErrMovieNotFound
	}

	return lastFilmID, nil
}

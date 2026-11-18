package movies

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/arinamklvch/go-courses/internal/data"
)

const workersCnt = 10

func GetRange(ctx context.Context, from int, to int, baseURL string, client *http.Client) ([]data.MovieResponse, int) {
	filmsCnt := to - from + 1
	jobs := make(chan int, filmsCnt)
	moviesRequestsResults := make(chan data.MovieRequestResult, filmsCnt)

	wg := sync.WaitGroup{}

	for range workersCnt {
		wg.Add(1)
		go worker(ctx, client, baseURL, jobs, moviesRequestsResults, &wg)
	}

jobsLoop:
	for filmID := from; filmID <= to; filmID++ {
		select {
		case <-ctx.Done():
			break jobsLoop
		case jobs <- filmID:
		}
	}
	close(jobs)

	wg.Wait()
	close(moviesRequestsResults)

	movieRangeResult := []data.MovieResponse{}

	successCount := 0
	for movieRequestResult := range moviesRequestsResults {
		if movieRequestResult.Error != nil {
			if errors.Is(movieRequestResult.Error, context.Canceled) {
				continue
			}

			log.Printf("load movies range (from=%d, to=%d): %v", from, to, movieRequestResult.Error)
			movieRangeResult = append(movieRangeResult, data.MovieResponse{
				ID:    movieRequestResult.ID,
				Error: movieRequestResult.Error.Error(),
			})

			continue
		}

		successCount++
		movieRangeResult = append(movieRangeResult, data.MovieResponse{
			ID:    movieRequestResult.ID,
			Movie: &movieRequestResult.Movie})
	}

	return movieRangeResult, successCount
}

func worker(ctx context.Context, client *http.Client, movieBaseURL string, jobs <-chan int, results chan<- data.MovieRequestResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		var filmID int
		var ok bool

		select {
		case <-ctx.Done():
			return
		case filmID, ok = <-jobs:
			if !ok {
				return
			}
		}

		movie, err := Download(ctx, client, movieBaseURL, filmID)
		downloadMovieResult := data.MovieRequestResult{
			ID:    filmID,
			Movie: movie,
			Error: err,
		}

		select {
		case <-ctx.Done():
			return
		case results <- downloadMovieResult:
		}
	}
}

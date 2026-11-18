package usecase

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/arinamklvch/go-courses/internal/domain"
)

const workersCnt = 10

func (s *Service) GetMoviesRange(ctx context.Context, from int, to int) ([]domain.MovieResult, int, error) {
	if from <= 0 {
		return nil, 0, domain.ErrIncorrectID
	}
	if from > to {
		return nil, 0, domain.ErrInvalidRange
	}

	filmsCnt := to - from + 1
	jobs := make(chan int, filmsCnt)
	moviesRequestsResponses := make(chan domain.MovieResult, filmsCnt)

	wg := sync.WaitGroup{}

	for range workersCnt {
		wg.Add(1)
		go worker(ctx, s.MovieProvider, jobs, moviesRequestsResponses, &wg)
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
	close(moviesRequestsResponses)

	movieRangeResult := []domain.MovieResult{}

	successCount := 0
	for movieRequestResponse := range moviesRequestsResponses {
		if movieRequestResponse.Error != nil {
			if errors.Is(movieRequestResponse.Error, context.Canceled) {
				continue
			}

			log.Printf("load movies range (from=%d, to=%d): %v", from, to, movieRequestResponse.Error)
		} else {
			successCount++
		}

		movieRangeResult = append(movieRangeResult, movieRequestResponse)

	}

	return movieRangeResult, successCount, nil
}

func worker(ctx context.Context, movieProvider MovieProvider, jobs <-chan int, responses chan<- domain.MovieResult, wg *sync.WaitGroup) {
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

		movie, err := movieProvider.GetMovie(ctx, filmID)
		downloadMovieResult := domain.MovieResult{
			ID:    filmID,
			Movie: movie,
			Error: err,
		}

		select {
		case <-ctx.Done():
			return
		case responses <- downloadMovieResult:
		}
	}
}

package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Movie struct {
	ID       int
	Title    string
	Year     int
	Director string
}

type Result struct {
	MovieID int
	Movie   Movie
	Err     error
}

const baseURL = "https://homeworksite.site/%d/info.0.json"

func validateFlags(from, to, workers, timeout int) error {
	seen := make(map[string]bool)

	flag.Visit(func(f *flag.Flag) {
		seen[f.Name] = true
	})

	for _, name := range []string{"from", "to"} {
		if !seen[name] {
			return fmt.Errorf("missing required flag --%s", name)
		}
	}

	if from > to {
		return fmt.Errorf("--from must be less than or equal to --to")
	}

	if workers <= 0 {
		return fmt.Errorf("--workers must be greater than 0")
	}

	if timeout <= 0 {
		return fmt.Errorf("--timeout must be greater than 0")
	}

	return nil
}

func downloadMovie(ctx context.Context, client *http.Client, filmID int) (Movie, error) {
	url := fmt.Sprintf(baseURL, filmID)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Movie{}, fmt.Errorf("create request: %w", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return Movie{}, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return Movie{}, fmt.Errorf("status %d", response.StatusCode)
	}

	var movie Movie
	if err := json.NewDecoder(response.Body).Decode(&movie); err != nil {
		return Movie{}, fmt.Errorf("decode JSON: %w", err)
	}

	return movie, nil
}

func worker(ctx context.Context, client *http.Client, jobs <-chan int, results chan<- Result, wg *sync.WaitGroup) {
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

		movie, err := downloadMovie(ctx, client, filmID)
		result := Result{
			MovieID: filmID,
			Movie:   movie,
			Err:     err,
		}

		select {
		case <-ctx.Done():
			return
		case results <- result:
		}
	}
}

func main() {
	from := flag.Int("from", 1, "left boundary")
	to := flag.Int("to", 100, "right boundary")
	workersCnt := flag.Int("workers", 10, "number of workers in worker pool")
	timeoutSeconds := flag.Int("timeout", 5, "timeout for HTTP request")

	flag.Parse()

	if err := validateFlags(*from, *to, *workersCnt, *timeoutSeconds); err != nil {
		log.Printf("flag error: %v", err)
		os.Exit(1)
	}

	filmsCnt := *to - *from + 1
	jobs := make(chan int, filmsCnt)
	results := make(chan Result, filmsCnt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	client := &http.Client{
		Timeout: time.Duration(*timeoutSeconds) * time.Second,
	}
	wg := sync.WaitGroup{}

	for range *workersCnt {
		wg.Add(1)
		go worker(ctx, client, jobs, results, &wg)
	}

	for filmID := *from; filmID <= *to; filmID++ {
		jobs <- filmID
	}
	close(jobs)

	wg.Wait()
	close(results)

	for result := range results {
		if result.Err != nil {
			if errors.Is(result.Err, context.Canceled) {
				continue
			}

			fmt.Fprintf(os.Stderr, "film %d: %v\n", result.MovieID, result.Err)
			continue
		}

		fmt.Printf("%d — %s — %d — %s\n", result.Movie.ID, result.Movie.Title, result.Movie.Year, result.Movie.Director)
	}

	if ctx.Err() != nil {
		fmt.Fprintln(os.Stderr, "interrupted: canceled remaining requests")
	}
}

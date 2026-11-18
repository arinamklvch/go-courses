package domain

type MovieResult struct {
	ID    int
	Movie Movie
	Error error
}

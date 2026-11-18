package movies

import "errors"

var ErrCreateMovieRequest = errors.New("create movie request")
var ErrMovieNotFound = errors.New("movie not found")

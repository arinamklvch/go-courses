package domain

import "errors"

var ErrCreateMovieRequest = errors.New("create movie request")
var ErrMovieNotFound = errors.New("movie not found")
var ErrInvalidRange = errors.New("from must be less than or equal to to")
var ErrIncorrectID = errors.New("from must be a positive integer")

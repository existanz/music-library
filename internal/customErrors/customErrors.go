package customErrors

import "errors"

var (
	ErrNotFound    = errors.New("song not found")
	ErrInvalidData = errors.New("invalid data")
	ErrISE         = errors.New("internal server error")
)

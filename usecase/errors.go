package usecase

import "errors"

// errors ...
var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidArguments = errors.New("invalid arguments")
)

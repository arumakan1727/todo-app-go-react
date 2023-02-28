package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")

	ErrNotFound = errors.New("requested item is not found")

	ErrConflict = errors.New("item already exist")

	ErrEmptyPatch = errors.New("no fields specified to patch")

	ErrInvalidInput = errors.New("given param is invalid")
)

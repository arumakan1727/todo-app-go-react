package domain

import "errors"

var (
	ErrNotFound = errors.New("requested item is not found")

	ErrAlreadyExit = errors.New("already exist")

	ErrEmptyPatch = errors.New("no fields specified to patch")

	ErrInvalidInput = errors.New("given param is invalid")

	ErrNotInTransaction = errors.New("not in transaction")
)

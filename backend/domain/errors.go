package domain

import "errors"

// アルファベット順
var (
	ErrAlreadyExits = errors.New("already exists")

	ErrEmptyPatch = errors.New("no fields specified to patch")

	ErrIncorrectEmailOrPasswd = errors.New("incorrect email or password")

	ErrInvalidInput = errors.New("given param is invalid")

	ErrNotFound = errors.New("requested item is not found")

	ErrNotInTransaction = errors.New("not in transaction")

	ErrUnauthorized = errors.New("unauthorized")
)

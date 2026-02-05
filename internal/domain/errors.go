package domain

import "errors"

var (
	ErrInternalProblems = errors.New("internal problems")
	ErrNotFound         = errors.New("Link not found")
)

package domain

import "errors"

var (
	ErrInternalProblems              = errors.New("internal problems")
	ErrNotFound                      = errors.New("Link not found")
	ErrLimitExceeded                 = errors.New("rate limit exceeded")
	ErrMaxAttemptsToGenerateShortUrl = errors.New("failed to generate unique short link after 6 attempts")
)

package bunq

import (
	"errors"
)

var (
	ErrInternalServerError        = errors.New("bunq: http request failed due to internal server error")
	ErrRateLimitExceeded          = errors.New("bunq: http request failed due to rate limit exceeded")
	ErrResponseVerificationFailed = errors.New("bunq: request was successful but response verification failed")
)

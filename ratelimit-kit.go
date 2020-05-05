package ratelimit_kit

import (
	"errors"
)

var (
	ErrExceededLimit = errors.New("Too many requests, exceeded the limit. ")
)

type RateLimiter interface {
	Take() error
}
